package digestwriter

// This source file contains an implementation of interface between Go code and
// (almost any) SQL database like PostgreSQL, SQLite, or MariaDB.

import (
	"app/base/models"
	"app/base/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultImageArch = "amd64"
)

// Storage represents an interface to almost any database or storage system
type Storage interface {
	WriteClusterInfo(cluster ClusterName, orgID AccountNumber, workload Workload, digests []string) error
}

// DBStorage is an implementation of Storage
// It is possible to configure connection via Configuration structure.
type DBStorage struct {
	connection *gorm.DB
	archMap    map[string]int64
}

func (storage *DBStorage) lookupArch(name string) error {
	archRows := []models.Arch{}
	if err := storage.connection.Where("name = ?", name).Find(&archRows).Error; err != nil {
		return err
	}
	for _, arch := range archRows {
		storage.archMap[arch.Name] = arch.ID
	}
	return nil
}

// NewStorage function creates and initializes a new instance of Storage interface
func NewStorage() (*DBStorage, error) {
	logger.Info("initializing connection to storage.")

	db, err := models.GetGormConnection(utils.GetDbURL(false))

	if err != nil {
		logger.Errorf("unable to connect to database: %s", err)
		return nil, err
	}

	logger.Infoln("connection to storage established")
	return NewFromConnection(db), nil
}

// NewFromConnection function returns a new Storage instance
// that will use the provided connection
func NewFromConnection(connection *gorm.DB) *DBStorage {
	return &DBStorage{
		connection: connection,
		archMap:    map[string]int64{},
	}
}

// prepareClusterImageLists recalculates previously inserted cluster images
// with newly obtained images and returns the differences
func prepareClusterImageLists(clusterID int64, currentImageIDs map[int64]struct{}, existingDigests []models.Image) (toInsert, toDelete []models.ClusterImage) {
	for _, digest := range existingDigests {
		if _, found := currentImageIDs[digest.ID]; !found {
			clusterImage := models.ClusterImage{ClusterID: clusterID, ImageID: digest.ID}
			toInsert = append(toInsert, clusterImage)
		} else {
			delete(currentImageIDs, digest.ID)
		}
	}
	for imageID := range currentImageIDs {
		clusterImage := models.ClusterImage{ClusterID: clusterID, ImageID: imageID}
		toDelete = append(toDelete, clusterImage)
	}
	return
}

// updateClusterCache updates the cache section of cluster row in db
func (storage *DBStorage) UpdateClusterCache(tx *gorm.DB, clusterID int64) error {
	subquery := tx.Table("cve").
		Select(`COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END), 0) AS c,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END), 0) AS i,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END), 0) AS m,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END), 0) AS l`,
			models.Critical, models.Important, models.Moderate, models.Low).
		Joins("JOIN image_cve ON image_cve.cve_id = cve.id").
		Joins("JOIN cluster_image ON cluster_image.image_id = image_cve.image_id").
		Where("cluster_image.cluster_id = ?", clusterID)

	res := tx.Exec(`UPDATE cluster SET cve_cache_critical = c.c, cve_cache_important = c.i, cve_cache_moderate = c.m, cve_cache_low = c.l FROM (?) AS c WHERE cluster.id = ?`, subquery, clusterID)
	if res.Error != nil {
		return fmt.Errorf("couldn't save cluster cache: %s", res.Error.Error())
	}

	return nil
}

// linkDigestsToCluster updates the 'cluster_image' table
func (storage *DBStorage) linkDigestsToCluster(tx *gorm.DB, clusterStr string, clusterID, clusterArchID int64, digests []string) error {
	//retrieve IDs of rows in image table for the received digests

	logger.Debugf("trying to link digests to cluster with ID %d", clusterID)

	var existingDigests []models.Image
	queryResult := tx.Where("(manifest_schema2_digest IN ? OR manifest_list_digest IN ? OR docker_image_digest IN ?) AND arch_id = ?",
		digests, digests, digests, clusterArchID).Find(&existingDigests)
	if err := queryResult.Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey:     err.Error(),
			clusterIDKey: clusterID,
		}).Errorln("couldn't retrieve any digest from table 'image' for the cluster with the given ID")
		return err
	}

	if queryResult.RowsAffected == 0 {
		logger.WithFields(logrus.Fields{
			clusterKey: clusterStr,
		}).Infoln("no digests in image table for the cluster with the given ID, nothing to do")
		return nil
	}

	logger.WithFields(logrus.Fields{
		clusterKey: clusterStr,
	}).Infof("linking %d digests from image table (%d/%d found)",
		queryResult.RowsAffected, queryResult.RowsAffected, len(digests),
	)

	var currentClusterImages []models.ClusterImage
	queryResult = tx.Where("cluster_id = ?", clusterID).Find(&currentClusterImages)
	if err := queryResult.Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey:     err.Error(),
			clusterIDKey: clusterID,
		}).Errorln("couldn't retrieve any rows from table 'cluster_image' for the cluster with the given ID")
		return err
	}
	currentImageIDs := map[int64]struct{}{}
	for _, clusterImage := range currentClusterImages {
		currentImageIDs[clusterImage.ImageID] = struct{}{}
	}
	toInsert, toDelete := prepareClusterImageLists(clusterID, currentImageIDs, existingDigests)

	if len(toInsert) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&toInsert).Error; err != nil {
			logger.WithFields(logrus.Fields{
				errorKey:    err.Error(),
				"clusterID": clusterID,
			}).Errorln("couldn't link cluster ID to image IDs for the given cluster")
			return err
		}
	}

	if len(toDelete) > 0 {
		deleteTx := tx.Session(&gorm.Session{})
		for _, clusterImage := range toDelete {
			deleteTx = deleteTx.Or(&clusterImage)
		}
		if err := deleteTx.Delete(&models.ClusterImage{}).Error; err != nil {
			logger.WithFields(logrus.Fields{
				errorKey:    err.Error(),
				"clusterID": clusterID,
			}).Errorln("couldn't unlink cluster ID from image IDs for the given cluster")
			return err
		}
	}

	err := storage.UpdateClusterCache(tx, clusterID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			errorKey:     err.Error(),
			clusterIDKey: clusterID,
		}).Errorln("couldn't update cluster cve cache")
		return err
	}

	logger.Debugln("linked digests to cluster successfully")
	return nil
}

// WriteClusterInfo updates the 'cluster' table with the provided info
func (storage *DBStorage) WriteClusterInfo(cluster ClusterName, orgID AccountNumber, workload Workload, digests []string) error {
	// prepare data
	clusterStr := string(cluster)
	clusterUUID, err := uuid.Parse(clusterStr)
	if err != nil {
		logger.Errorln("cannot convert given cluster ID to UUID, aborting WriteClusterInfo")
		return err
	}
	accountData := models.Account{
		OrgID: string(orgID),
	}

	logger.WithFields(logrus.Fields{
		orgKey: accountData.OrgID,
	}).Debugln("account data to insert")

	tx := storage.connection.Begin()

	// Insert account info in account table if not present
	// If present, retrieve corresponding ID
	if err = tx.Where(accountData).
		Clauses(clause.OnConflict{DoNothing: true}).
		FirstOrCreate(&accountData).Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err.Error(),
		}).Errorln("couldn't insert or retrieve cluster name in 'account' table")
		if r := tx.Rollback(); r.Error != nil {
			logger.WithFields(logrus.Fields{
				errorKey: r.Error.Error(),
			}).Errorln("couldn't rollback operation!")
			return r.Error
		}
		return err
	}

	logger.WithFields(logrus.Fields{
		rowIDKey: accountData.ID,
		orgKey:   accountData.OrgID,
	}).Debugln("inserted account data successfully")

	clusterInfoData := models.Cluster{
		UUID:      clusterUUID,
		AccountID: accountData.ID,
		LastSeen:  time.Now().UTC(),
	}

	if err := clusterInfoData.Workload.Set(workload); err != nil {
		logger.Errorln("cannot set workload JSON")
		return err
	}

	if err := tx.Omit(
		"DisplayName", "Status", "Type", "Version", "Provider", "CveCacheCritical",
		"CveCacheImportant", "CveCacheModerate", "CveCacheLow").
		Where(clusterInfoData).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uuid"}},
			UpdateAll: true,
		}).FirstOrCreate(&clusterInfoData).Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err.Error(),
		}).Errorln("couldn't write cluster info in cluster table")
		return err
	}

	logger.WithFields(logrus.Fields{
		rowIDKey:   clusterInfoData.ID,
		clusterKey: clusterInfoData.UUID,
		accountKey: clusterInfoData.AccountID,
	}).Debugln("updated cluster table successfully")

	var archID int64
	var found bool
	if archID, found = storage.archMap[defaultImageArch]; !found {
		// TODO: get real cluster image arch, use default for now
		err = storage.lookupArch(defaultImageArch)
		if err != nil {
			return err
		}
		if archID, found = storage.archMap[defaultImageArch]; !found {
			err = fmt.Errorf("unknown image arch: %s", defaultImageArch)
			return err
		}
	}

	if err = storage.linkDigestsToCluster(tx, clusterStr, clusterInfoData.ID, archID, digests); err != nil {
		return tx.Rollback().Error
	}
	return tx.Commit().Error
}

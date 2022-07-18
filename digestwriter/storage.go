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

// Storage represents an interface to almost any database or storage system
type Storage interface {
	WriteClusterInfo(cluster ClusterName, account AccountNumber, orgID OrgID, digests []string) error
}

// DBStorage is an implementation of Storage
// It is possible to configure connection via Configuration structure.
type DBStorage struct {
	connection *gorm.DB
}

// NewStorage function creates and initializes a new instance of Storage interface
func NewStorage() (*DBStorage, error) {
	logger.Info("Initializing connection to storage.")

	db, err := models.GetGormConnection(utils.GetDbURL(false))

	if err != nil {
		logger.Errorf("Unable to connect to database: %s\n", err)
		return nil, err
	}

	logger.Infoln("Connection to storage established")
	return NewFromConnection(db), nil
}

// NewFromConnection function returns a new Storage instance
// that will use the provided connection
func NewFromConnection(connection *gorm.DB) *DBStorage {
	return &DBStorage{
		connection: connection,
	}
}

func prepareBulkInsertClusterImage(clusterID int64, digests []models.Image) (data []models.ClusterImage) {
	data = make([]models.ClusterImage, len(digests))
	for idx, digest := range digests {
		data[idx].ClusterID = clusterID
		data[idx].ImageID = digest.ID
	}
	return
}

// linkDigestsToCluster updates the 'cluster_image' table
func (storage *DBStorage) linkDigestsToCluster(tx *gorm.DB, clusterID int64, digests []string) error {
	//retrieve IDs of rows in image table for the received digests

	logger.Infof("Trying to link digests to cluster with ID %d\n", clusterID)

	var existingDigests []models.Image
	queryResult := tx.Where("digest IN ?", digests).Find(&existingDigests)
	if err := queryResult.Error; err != nil {
		//TODO: Maybe we prefer to check digests first, and not insert anything in cluster and cluster_image tables?
		logger.WithFields(logrus.Fields{
			errorKey:     err.Error(),
			clusterIDKey: clusterID,
		}).Errorln("couldn't retrieve any digest from table 'image' for the cluster with the given ID")
		return err
	}

	if queryResult.RowsAffected == 0 {
		logger.WithFields(logrus.Fields{
			clusterIDKey: clusterID,
		}).Infoln("no digests in image table for the cluster with the given ID. Nothing to do.")
		return nil
	}

	logger.Infof("Found %d digests in image table (%d/%d)\n",
		queryResult.RowsAffected, queryResult.RowsAffected, len(digests),
	)

	clusterImageData := prepareBulkInsertClusterImage(clusterID, existingDigests)

	// Do nothing on conflict. It just means that we already have
	// the info we are trying to insert
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&clusterImageData).Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey:    err.Error(),
			"clusterID": clusterID,
		}).Errorln("couldn't link cluster ID to image IDs for the given cluster")
		return err
	}

	logger.Infoln("linked digests to cluster successfully")
	return nil
}

// WriteClusterInfo updates the 'cluster' table with the provided info
func (storage *DBStorage) WriteClusterInfo(cluster ClusterName, account AccountNumber, orgID OrgID, digests []string) error {
	// prepare data
	clusterUUID, err := uuid.Parse(string(cluster))
	if err != nil {
		logger.Errorln("cannot convert given cluster ID to UUID. Aborting WriteClusterInfo")
		return err
	}
	accountData := models.Account{
		AccountNumber: fmt.Sprint(account),
		OrgID:         fmt.Sprint(orgID),
	}

	logger.WithFields(logrus.Fields{
		rowIDKey:   accountData.ID,
		accountKey: accountData.AccountNumber,
		orgKey:     accountData.OrgID,
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
		rowIDKey:   accountData.ID,
		accountKey: accountData.AccountNumber,
		orgKey:     accountData.OrgID,
	}).Debugln("inserted account data successfully")

	clusterInfoData := models.Cluster{
		UUID:      clusterUUID,
		AccountID: accountData.ID,
		LastSeen:  time.Now().UTC(),
	}

	if err := tx.Omit(
		"Status", "Version", "Provider", "CveCacheCritical",
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

	if err = storage.linkDigestsToCluster(tx, clusterInfoData.ID, digests); err != nil {
		return tx.Rollback().Error
	}
	return tx.Commit().Error
}

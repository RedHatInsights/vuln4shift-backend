package vmsync

import (
	"app/base/logging"
	"app/base/models"
	"app/base/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	logger    *logrus.Logger
	BatchSize = utils.GetIntEnv("BATCH_SIZE", 5000)
)

func init() {
	logLevel := utils.GetEnv("LOGGING_LEVEL", "INFO")
	var err error
	logger, err = logging.CreateLogger(logLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func syncCveMetadata() {
	apiCveMap, err := getAPICves()

	if err != nil {
		logger.Fatalf("Unable to get CVEs from VMaaS: %s", err)
	}
	logger.Infof("CVEs in VMaaS: %d", len(apiCveMap))

	toSyncCves := make([]models.Cve, 0, len(apiCveMap))
	for cveName, apiCve := range apiCveMap {
		var severity models.Severity
		err := severity.Scan(apiCve.Impact)
		if err != nil {
			severity = models.NotSet
		}

		cvss2Score, _ := strconv.ParseFloat(apiCve.Cvss2Score, 32)
		cvss3Score, _ := strconv.ParseFloat(apiCve.Cvss3Score, 32)

		toSyncCves = append(toSyncCves, models.Cve{
			Name:         cveName,
			Description:  apiCve.Description,
			Severity:     severity,
			Cvss2Metrics: apiCve.Cvss2Metrics,
			Cvss2Score:   float32(cvss2Score),
			Cvss3Metrics: apiCve.Cvss3Metrics,
			Cvss3Score:   float32(cvss3Score),
			PublicDate:   apiCve.PublicDate,
			ModifiedDate: apiCve.ModifiedDate,
			RedhatURL:    apiCve.RedhatURL,
			SecondaryURL: apiCve.SecondaryURL,
		})

		if _, found := dbCveMap[cveName]; found {
			delete(dbCveMap, cveName)
		}
	}

	toDeleteCves := make([]models.Cve, 0, len(dbCveMap))
	for _, dbCve := range dbCveMap {
		toDeleteCves = append(toDeleteCves, dbCve)
	}

	logger.Infof("CVEs to sync: %d", len(toSyncCves))
	logger.Infof("CVEs to delete: %d", len(toDeleteCves))

	if err = syncCves(toSyncCves, toDeleteCves); err != nil {
		logger.Infof("Error during syncing CVEs into database: %s", err)
	}

	logger.Infof("Metadata sync finished successfully")
}

func syncCves(toSyncCves, toDeleteCves []models.Cve) error {
	tx := DB.Begin()
	// Do a rollback by default (don't need to specify on every return), will do nothing when everything is committed
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if len(toSyncCves) > 0 {
		if err := insertUpdateCves(toSyncCves, tx); err != nil {
			return errors.Wrap(err, "Unable to insert/update cves in database")
		}
	}

	if len(toDeleteCves) > 0 {
		if err := deleteCves(toDeleteCves, tx); err != nil {
			return errors.Wrap(err, "Unable to delete cves in database")
		}
	}

	return tx.Commit().Error
}

func deleteCves(toDeleteCves []models.Cve, tx *gorm.DB) error {
	logger.Debugf("CVEs to delete: %d", len(toDeleteCves))

	ids := make([]int64, 0, len(toDeleteCves))
	for _, cve := range toDeleteCves {
		ids = append(ids, cve.ID)
	}

	if err := tx.Where("cve_id in ?", ids).Delete(&models.AccountCveCache{}).Error; err != nil {
		return err
	}

	if err := tx.Where("cve_id in ?", ids).Delete(&models.ClusterCveCache{}).Error; err != nil {
		return err
	}

	if err := tx.Where("cve_id in ?", ids).Delete(&models.ImageCve{}).Error; err != nil {
		return err
	}

	if err := tx.Delete(&toDeleteCves).Error; err != nil {
		return err
	}
	return nil
}

func insertUpdateCves(toSyncCves []models.Cve, tx *gorm.DB) error {
	logger.Debugf("CVEs to insert/update: %d", len(toSyncCves))

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).CreateInBatches(toSyncCves, BatchSize).Error; err != nil {
		return err
	}
	return nil
}

func Start() {
	logger.Info("Starting vmaas sync.")

	if err := dbConfigure(); err != nil {
		logger.Fatalf("Unable to get GORM connection: %s", err)
	}
	if err := prepareDbCvesMap(); err != nil {
		logger.Fatalf("Unable to fetch data from DB: %s", err)
	}

	logger.Infof("CVEs in DB: %d", len(dbCveMap))

	syncCveMetadata()
}

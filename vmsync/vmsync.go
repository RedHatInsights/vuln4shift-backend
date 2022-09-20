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
)

var (
	logger    *logrus.Logger
	BatchSize = utils.Cfg.VmaasBatchSize
)

func init() {
	var err error
	logger, err = logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func syncCveMetadata() {
	apiCveSortedList, apiCveMap, err := getAPICves()

	if err != nil {
		logger.Fatalf("Unable to get CVEs from VMaaS: %s", err)
	}

	toSyncCves := make([]models.Cve, 0, len(apiCveMap))
	for _, cveName := range apiCveSortedList {
		apiCve := apiCveMap[cveName]
		var severity models.Severity
		err := severity.Scan(apiCve.Impact)
		if err != nil {
			severity = models.NotSet
		}

		cvss2Score, err := strconv.ParseFloat(apiCve.Cvss2Score, 32)
		if err != nil {
			cvss2Score = 0.0
		}
		cvss3Score, err := strconv.ParseFloat(apiCve.Cvss3Score, 32)
		if err != nil {
			cvss3Score = 0.0
		}

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

	logger.Infof("CVEs to sync: %d", len(toSyncCves))

	if err = syncCves(toSyncCves); err != nil {
		logger.Fatalf("Error during syncing CVEs into database: %s", err)
	}

	if err = pruneCves(); err != nil {
		logger.Fatalf("Failed to prune CVEs: %s", err)
	}

	logger.Infof("Metadata sync finished successfully")
}

// syncCves inserts or updates CVEs into the database.
func syncCves(toSyncCves []models.Cve) error {
	tx := DB.Begin()
	// Do a rollback by default (don't need to specify on every return), will do nothing when everything is committed
	defer tx.Rollback()

	toSyncCvesCnt := len(toSyncCves)
	if toSyncCvesCnt > 0 {
		if err := insertUpdateCves(toSyncCves, tx); err != nil {
			syncError.WithLabelValues(dbInsertUpdate).Inc()
			return errors.Wrap(err, "Unable to insert/update cves in database")
		}
		cvesInsertedUpdated.Add(float64(toSyncCvesCnt))
	}

	return tx.Commit().Error
}

// pruneCves deletes CVEs from DB absent in recent VMaaS response.
func pruneCves() error {
	tx := DB.Begin()
	defer tx.Rollback()

	notInVmaasCves := make([]int64, 0, len(dbCveMap))
	for _, dbCve := range dbCveMap {
		notInVmaasCves = append(notInVmaasCves, dbCve.ID)
	}

	if len(notInVmaasCves) > 0 {
		var deletedCnt int64
		var err error
		if deletedCnt, err = deleteNotAffectingCves(tx, notInVmaasCves); err != nil {
			syncError.WithLabelValues(dbDelete).Inc()
			return errors.Wrap(err, "failed to delete CVEs from database")
		}
		cvesDeleted.Add(float64(deletedCnt))
		logger.Infof("Deleted %d CVEs from database", deletedCnt)
	}

	return tx.Commit().Error
}

func Start() {
	logger.Info("Starting vmaas sync.")

	pusher := GetMetricsPusher()

	if err := dbConfigure(); err != nil {
		syncError.WithLabelValues(dbConnection).Inc()
		logger.Fatalf("Unable to get GORM connection: %s", err)
	}
	if err := prepareDbCvesMap(); err != nil {
		syncError.WithLabelValues(dbFetch).Inc()
		logger.Fatalf("Unable to fetch data from DB: %s", err)
	}

	logger.Infof("CVEs in DB: %d", len(dbCveMap))

	syncCveMetadata()

	if err := pusher.Add(); err != nil {
		logger.Warningln("Could not push to Pushgateway:", err)
	}
}

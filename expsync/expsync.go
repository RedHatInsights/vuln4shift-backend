package expsync

import (
	"app/base/logging"
	"app/base/models"
	"app/base/utils"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	var err error
	logger, err = logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		logger.Fatalf("Error setting up logger.")
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// pruneCvesMetadata removes exploit metadata from CVEs absent in the API response.
func pruneCvesMetadata(tx *gorm.DB, apiExploits map[CVE][]ExploitMetadata) error {
	dbCvesWithExploit, err := getCvesWithExploitMetadata(tx)
	if err != nil {
		return errors.Wrap(err, "failed to get CVEs with exploit metadata")
	}

	var cvesToPrune []models.Cve

	for _, dbCve := range dbCvesWithExploit {
		if _, found := apiExploits[CVE(dbCve.Name)]; !found {
			cvesToPrune = append(cvesToPrune, dbCve)
		}
	}

	removedCnt, err := removeExploitData(tx, cvesToPrune)
	if err != nil {
		return err
	}

	if removedCnt > 0 {
		cveExploitsDeleted.Add(float64(removedCnt))
		logger.Infof("Pruned exploits from DB: %d", removedCnt)
	}

	return nil
}

// syncExploits updates DB with exploit metadata from remote source.
func syncExploits() error {
	tx := DB.Begin()
	defer tx.Rollback()

	exploitsMetadata, err := getAPIExploits()
	if err != nil {
		return errors.Wrap(err, "unable to get exploit metadata file")
	}
	logger.Infof("Exploits fetched from remote source: %d", len(exploitsMetadata))

	err = pruneCvesMetadata(tx, exploitsMetadata)
	if err != nil {
		syncError.WithLabelValues(dbDelete).Inc()
		return errors.Wrap(err, "failed to prune CVEs exploit metadata")
	}

	updatedCnt, err := updateExploitsMetadata(tx, exploitsMetadata)
	if err != nil {
		syncError.WithLabelValues(dbInsertUpdate).Inc()
		return errors.Wrap(err, "failed to update exploits metadata in DB")
	}
	cveExploitsInsertedUpdated.Add(float64(updatedCnt))
	logger.Infof("Exploits updated: %d", updatedCnt)

	if diff := int64(len(exploitsMetadata)) - updatedCnt; diff > 0 {
		logger.Infof("Missing CVEs in DB versus source: %d", diff)
	}

	return tx.Commit().Error
}

func Start() {
	logger.Info("Starting exploit sync.")

	var err error
	DB, err = utils.DbConfigure()
	if err != nil {
		syncError.WithLabelValues(dbConnection).Inc()
		logger.Fatalf("Unable to get GORM connection: %s", err)
	}

	err = syncExploits()
	if err != nil {
		logger.Fatalf("exploit sync failed: %s", err)
	}

	pusher := getMetricsPusher()

	if err := pusher.Add(); err != nil {
		logger.Warningln("Could not push to Pushgateway:", err)
	}
}

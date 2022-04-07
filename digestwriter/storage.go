package digestwriter

// This source file contains an implementation of interface between Go code and
// (almost any) SQL database like PostgreSQL, SQLite, or MariaDB.

import (
	"app/base/models"
	"app/base/utils"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Storage represents an interface to almost any database or storage system
type Storage interface {
	WriteDigests(digests []string) error
}

// DBStorage is an implementation of Storage
// It is possible to configure connection via Configuration structure.
type DBStorage struct {
	connection *gorm.DB
}

// NewStorage function creates and initializes a new instance of Storage interface
func NewStorage() (*DBStorage, error) {
	logger.Info("Initializing connection to storage.")

	db, err := models.GetGormConnection(utils.GetDbURL())

	if err != nil {
		logger.Errorf("Unable to connect to database: %s\n", err)
		return nil, err
	}

	logger.Infoln("Connection to storage established")
	return NewFromConnection(db), nil
}

// NewFromConnection function creates and initializes a new instance of Storage interface from prepared connection
func NewFromConnection(connection *gorm.DB) *DBStorage {
	return &DBStorage{
		connection: connection,
	}
}

func prepareBulkInsertDigestsStruct(digests []string) (data []models.Image) {
	data = make([]models.Image, len(digests))
	modifiedAt := time.Now().UTC()
	for idx, digest := range digests {
		data[idx].Digest = digest
		data[idx].ModifiedDate = modifiedAt
	}
	return
}

// WriteDigests writes digests into the 'image' table
func (storage *DBStorage) WriteDigests(digests []string) error {
	data := prepareBulkInsertDigestsStruct(digests)

	logger.WithFields(logrus.Fields{
		"num_rows": len(data),
	}).Debug("trying to insert digests.")

	// Begin a new transaction.
	tx := storage.connection.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Omit("PyxisID").Create(&data).Error; err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err,
		}).Debug("Couldn't insert digests.")
		return err
	}

	return tx.Commit().Error
}

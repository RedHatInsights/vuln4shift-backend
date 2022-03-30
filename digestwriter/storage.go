package digestwriter

// This source file contains an implementation of interface between Go code and
// (almost any) SQL database like PostgreSQL, SQLite, or MariaDB.

import (
	"app/base/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"app/base/utils"
)

// Storage represents an interface to almost any database or storage system
type Storage interface {
	WriteDigests(digests []string) error
}

// DBStorage is an implementation of Storage
// It is possible to configure connection via Configuration structure.
type DBStorage struct {
	connection   *gorm.DB
	Logger       *logrus.Logger
}

// NewStorage function creates and initializes a new instance of Storage interface
func NewStorage(logger *logrus.Logger) (*DBStorage, error) {
	logger.Info("Initializing connection to storage.")

		dsn := utils.GetDbURL()
		db, err := models.GetGormConnection(dsn)

	if err != nil {
		logger.Fatalf("Unable to connect to database: %s.\n", err)
		return nil, err
	}

	logger.Info("Connection to storage established.")
	return NewFromConnection(db, logger), nil
}

// NewFromConnection function creates and initializes a new instance of Storage interface from prepared connection
func NewFromConnection(connection *gorm.DB, logger *logrus.Logger) *DBStorage {
	return &DBStorage{
		connection:   connection,
		Logger: 	  logger,
	}
}

func prepareBulkInsertDigestsStruct(digests []string) (data []models.Image) {
	data = make([]models.Image, len(digests))
	for idx, digest := range digests {
		data[idx].Digest = digest
	}
	return
}

// WriteDigests writes digests into the 'image' table
func (storage *DBStorage) WriteDigests(digests []string) error {

	data := prepareBulkInsertDigestsStruct(digests)

	storage.Logger.WithFields(logrus.Fields{
			"num_rows": len(data),
		}).Debug("trying to insert digests.")

	// Begin a new transaction.
	tx := storage.connection.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Omit("HealthIndex").Create(&data).Error; err != nil {
		storage.Logger.WithFields(logrus.Fields{
			errorKey: err,
		}).Debug("Couldn't insert digests.")
		return err
	}

	return tx.Commit().Error
}


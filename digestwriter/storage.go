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
	Close() error
	WriteDigests(
		digests []string,
	) error
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

// Close method closes the connection to database.
// This only works with a single master connection, but when dealing with
// replicas using DBResolver, it does not close everything since gorm.DB.DB()
// only returns the master connection.
func (storage DBStorage) Close() error {
	storage.Logger.Info("Closing connection to data storage.")
	db, err := storage.connection.DB()
	if err != nil {
		storage.Logger.Fatalf("Unable to retrieve connection to data storage: %s\n", err)
		return err
	}
	err = db.Close()
	if err != nil {
		storage.Logger.Fatalf("Cannot close connection to data storage: %s\n", err)
		return err
	}
	return nil
}

func prepareBulkInsertDigestsStruct(digests []string) (data []models.Image) {
	data = make([]models.Image, len(digests))
	for idx, digest := range digests {
		data[idx].Digest = digest
	}
	return
}

// WriteDigests writes digests into the 'image' table
func (storage DBStorage) WriteDigests(digests []string) error {

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


package digestwriter_test

// Unit test definitions for functions and methods defined in source file
// storage.go

import (
	"app/base/logging"
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"app/digestwriter"
)

// mustCreateMockConnection function tries to create a new mock connection and
// checks if the operation was finished without problems.
func mustCreateMockConnection(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	// try to initialize new mock connection
	connection, mock, err := sqlmock.New()

	// check the status
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return connection, mock
}

func createGormMockConnection(t *testing.T, db *sql.DB) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	return gorm.Open(dialector, &gorm.Config{})
}

// checkConnectionClose function perform mocked DB closing operation and checks
// if the connection is properly closed from unit tests.
func checkConnectionClose(t *testing.T, connection *sql.DB) {
	// connection to mocked DB needs to be closed properly
	err := connection.Close()

	// check the error status
	if err != nil {
		t.Fatalf("error during closing connection: %v", err)
	}
}

// checkAllExpectations function checks if all database-related operations have
// been really met.
func checkAllExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	// check if all expectations were met
	err := mock.ExpectationsWereMet()

	// check the error status
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestWriteSingleDigest function tests the method Storage.WriteDigests
// when only one digest is passed in the slice
func TestWriteSingleDigest(t *testing.T) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)
	DB, err := createGormMockConnection(t, connection)
	if err != nil {
		t.Errorf("error was not expected while creating mock connection: %s", err)
	}

	logger, err := logging.CreateLogger("DEBUG")
	// prepare connection to mocked database
	storage := digestwriter.NewFromConnection(DB, logger)

	// expected SQL statements during this test
	//expectedStatement := `INSERT INTO "image" ("digest") VALUES ($1);`
	expectedStatement := `INSERT INTO "image" ("digest") VALUES ($1) RETURNING "id"`

	mock.ExpectBegin()
	//mock.ExpectExec(expectedStatement).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs("digest1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	mock.ExpectClose()

	// call the tested method
	err = storage.WriteDigests([]string{"digest1"})
	if err != nil {
		t.Errorf("error was not expected while writing the digests: %s", err)
	}

	// connection to mocked DB needs to be closed properly
	checkConnectionClose(t, connection)

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

// TestWriteSingleDigest function tests the method Storage.WriteDigests
// when multiple digests are passed in the slice
func TestWriteMultipleDigest(t *testing.T) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)
	DB, err := createGormMockConnection(t, connection)
	if err != nil {
		t.Errorf("error was not expected while creating mock connection: %s", err)
	}

	logger, err := logging.CreateLogger("DEBUG")
	// prepare connection to mocked database
	storage := digestwriter.NewFromConnection(DB, logger)

	// expected SQL statements during this test
	expectedStatement := `INSERT INTO "image" ("digest") VALUES ($1),($2) RETURNING "id"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs("digest1", "digest2").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
	mock.ExpectCommit()
	mock.ExpectClose()

	// call the tested method
	err = storage.WriteDigests([]string{"digest1", "digest2"})
	if err != nil {
		t.Errorf("error was not expected while writing the digests: %s", err)
	}

	// connection to mocked DB needs to be closed properly
	checkConnectionClose(t, connection)

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

// TestClose function tests database close operation.
func TestClose(t *testing.T) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)
	DB, err := createGormMockConnection(t, connection)
	if err != nil {
		t.Errorf("error was not expected while creating mock connection: %s", err)
	}

	logger, err := logging.CreateLogger("DEBUG")
	// prepare connection to mocked database
	storage := digestwriter.NewFromConnection(DB, logger)

	// we just happen to make connection without trying to actually connect
	assert.Nil(t, err)

	// try to close the storage
	mock.ExpectClose()
	err = storage.Close()

	// it should not fail
	assert.Nil(t, err)
}

package digestwriter_test

// Unit test definitions for functions and methods defined in source file
// storage.go

import (
	"app/digestwriter"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// mustCreateMockConnection function tries to create a new mock connection and
// checks if the operation was finished without problems.
func mustCreateMockConnection(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	//initialize logger
	digestwriter.SetupLogger()
	// try to initialize new mock connection
	connection, mock, err := sqlmock.New()

	// check the status
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return connection, mock
}

func createGormMockPostgresConnection(db *sql.DB) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	return gorm.Open(dialector, &gorm.Config{})
}

func NewMockStorage(t *testing.T) (*digestwriter.DBStorage, sqlmock.Sqlmock) {
	// prepare new mocked connection to database
	connection, mock := mustCreateMockConnection(t)
	DB, err := createGormMockPostgresConnection(connection)
	if err != nil {
		t.Errorf("error was not expected while creating mock connection: %s", err)
	}

	// prepare connection to mocked database
	return digestwriter.NewFromConnection(DB), mock
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
	storage, mock := NewMockStorage(t)
	patchCurrentTime()
	// expected SQL statements during this test
	expectedStatement := `INSERT INTO "image" ("modified_date","digest") VALUES ($1,$2) RETURNING "id"`

	mock.ExpectBegin()
	//mock.ExpectExec(expectedStatement).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs(time.Now().UTC(), "digest1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteDigests([]string{"digest1"})
	if err != nil {
		t.Errorf("error was not expected while writing the digests: %s", err)
	}

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

// TestWriteSingleDigest function tests the method Storage.WriteDigests
// when multiple digests are passed in the slice
func TestWriteMultipleDigest(t *testing.T) {
	storage, mock := NewMockStorage(t)
	patchCurrentTime()
	// expected SQL statements during this test
	expectedStatement := `INSERT INTO "image" ("modified_date","digest") VALUES ($1,$2),($3,$4) RETURNING "id"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs(time.Now().UTC(), "digest1", time.Now().UTC(), "digest2").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteDigests([]string{"digest1", "digest2"})
	if err != nil {
		t.Errorf("error was not expected while writing the digests: %s", err)
	}

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

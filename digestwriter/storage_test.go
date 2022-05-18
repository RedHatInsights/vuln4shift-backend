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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	digestModifiedAtTime = time.Now().UTC()
	clusterName          = digestwriter.ClusterName("84f7eedc-0000-0000-9d4d-000000000000")
	invalidClusterName   = digestwriter.ClusterName("99z7zzzz-0000-0000-9d4d-000000000000")
	testAccountNumber    = digestwriter.AccountNumber(1)
	testOrgID            = digestwriter.OrgID(1)

	anyArgForMockSQLQueries = sqlmock.AnyArg()
)

const (
	testClustedID = 1
	testAccountID = 1

	firstDigestID  = 1
	firstDigest    = "digest1"
	secondDigestID = 3
	secondDigest   = "digest2"
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

// TestLinkSingleDigestToCluster function tests the method Storage.linkDigestsToCluster
// when only one digest is passed in the slice
func TestLinkSingleDigestToCluster(t *testing.T) {
	storage, mock := NewMockStorage(t)
	// expected SQL statements during this test
	expectedSelect := `SELECT * FROM "image" WHERE digest IN ($1)`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest).
		WillReturnRows(sqlmock.NewRows([]string{"id", "digest", "pyxis_id", "modified_date"}).
			AddRow(firstDigestID, firstDigestID, 5, digestModifiedAtTime))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(firstDigestID, 1))
	mock.ExpectCommit()

	// call the tested method
	err := digestwriter.LinkDigestsToCluster(storage, testClustedID, []string{firstDigest})
	if err != nil {
		t.Errorf("error was not expected while linking the digests: %s", err)
	}

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

// TestLinkMultipleDigestsToCluster function tests the method Storage.linkDigestToCluster
// when multiple digests are passed in the slice
func TestLinkMultipleDigestsToCluster(t *testing.T) {
	storage, mock := NewMockStorage(t)
	// expected SQL statements during this test
	expectedSelect := `SELECT * FROM "image" WHERE digest IN ($1,$2)`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest, secondDigest).
		WillReturnRows(sqlmock.NewRows([]string{"id", "digest", "pyxis_id", "modified_date"}).
			AddRow(firstDigestID, firstDigest, 5, digestModifiedAtTime).
			AddRow(secondDigestID, secondDigest, 6, digestModifiedAtTime),
		)
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(secondDigestID, 2))
	mock.ExpectCommit()

	// call the tested method
	err := digestwriter.LinkDigestsToCluster(storage, 1, []string{firstDigest, secondDigest})
	if err != nil {
		t.Errorf("error was not expected while linking the digests: %s", err)
	}

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfo function tests that the method Storage.WriteClusterInfo
// does not modify the DB if the received clusterName is not valid
func TestWriteClusterInfoWrongClusterName(t *testing.T) {
	storage, mock := NewMockStorage(t)

	// call the tested method
	err := storage.WriteClusterInfo(&invalidClusterName, &testAccountNumber, &testOrgID, []string{firstDigest})
	assert.Error(t, err, "the given UUID should not have been parsed")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfoWithExistingAccountForClusterName tests happy path for
// Storage.WriteClusterInfo (account found for given cluster)
func TestWriteClusterInfoWithExistingAccountForClusterName(t *testing.T) {
	storage, mock := NewMockStorage(t)
	mock.ExpectBegin()

	//Expect the select query and return 1 records
	expectedSelect := `SELECT * FROM "account" WHERE "account"."account_number" = $1 AND "account"."org_id" = $2 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testAccountID))

	// Expect a SELECT FROM cluster statement since the previous select returned a record
	// Since it returns a record, there is no need to expect the INSERT INTO "cluster" statement
	expectedSelect = `SELECT "cluster"."id","cluster"."uuid","cluster"."account_id" FROM "cluster" WHERE "cluster"."uuid" = $1 AND "cluster"."account_id" = $2 ORDER BY "cluster"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(clusterName, testAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "account_id"}).AddRow(testClustedID, clusterName, testAccountID))

	// Since it all went smoothly, the digest is linked to the cluster
	expectedSelect = `SELECT * FROM "image" WHERE digest IN ($1)`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest).
		WillReturnRows(sqlmock.NewRows([]string{"id", "digest", "pyxis_id", "modified_date"}).
			AddRow(firstDigestID, firstDigestID, 5, digestModifiedAtTime))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(firstDigestID, 1))

	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteClusterInfo(&clusterName, &testAccountNumber, &testOrgID, []string{firstDigest})
	assert.Nil(t, err, "No error expected")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfoNoAccountForClusterName tests that
// Storage.WriteClusterInfo creates a new record in account table if
//  no account can be retrieved for given clusterName and updates
//  the cluster_image table properly
func TestWriteClusterInfoNoAccountForClusterName(t *testing.T) {
	storedAccountID := 10
	storedClusterID := 2

	storage, mock := NewMockStorage(t)
	mock.ExpectBegin()

	//Expect the select query and mock that it returns 0 records
	expectedSelect := `SELECT * FROM "account" WHERE "account"."account_number" = $1 AND "account"."org_id" = $2 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	//Expect an 'INSERT INTO account' statement since the select returned no rows
	expectedInsert := `INSERT INTO "account" ("account_number","org_id") VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING "id"`
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsert)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storedAccountID))

	// Expect a SELECT FROM cluster statement
	// Mock that it would not return any row since there was no data for current account
	expectedSelect = `SELECT "cluster"."id","cluster"."uuid","cluster"."account_id" FROM "cluster" WHERE "cluster"."uuid" = $1 AND "cluster"."account_id" = $2 ORDER BY "cluster"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "account_id"}))
	//Expect an 'INSERT INTO cluster' with the ID of the created account record
	expectedInsert = `INSERT INTO "cluster" ("uuid","account_id") VALUES ($1,$2) ON CONFLICT DO NOTHING RETURNING "id"`
	clusterUUID, _ := uuid.Parse(string(clusterName))
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsert)).WithArgs(clusterUUID, storedAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storedClusterID))

	// Since it all went smoothly, the digest is linked to the cluster
	expectedSelect = `SELECT * FROM "image" WHERE digest IN ($1)`
	expectedInsert = `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest).
		WillReturnRows(sqlmock.NewRows([]string{"id", "digest", "pyxis_id", "modified_date"}).
			AddRow(firstDigestID, firstDigestID, 5, digestModifiedAtTime))

	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WithArgs(storedClusterID, firstDigestID).
		WillReturnResult(sqlmock.NewResult(firstDigestID, 1))

	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteClusterInfo(&clusterName, &testAccountNumber, &testOrgID, []string{firstDigest})
	assert.Nil(t, err, "No error expected")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfoErrorWritingAccount tests that
// Storage.WriteClusterInfo aborts and DB is rolled back to
// previous state if an error happens while getting or creating
//  the account info
func TestWriteClusterInfoErrorWritingAccount(t *testing.T) {
	storage, mock := NewMockStorage(t)
	mock.ExpectBegin()

	// expect select from account statement
	expectedSelect := `SELECT * FROM "account" WHERE "account"."account_number" = $1 AND "account"."org_id" = $2 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnError(sql.ErrConnDone)

	mock.ExpectRollback()

	// call the tested method
	err := storage.WriteClusterInfo(&clusterName, &testAccountNumber, &testOrgID, []string{firstDigest})
	assert.Error(t, err, "account table shouldn't be updated and WriteClusterInfo should abort")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

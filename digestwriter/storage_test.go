package digestwriter_test

// Unit test definitions for functions and methods defined in source file
// storage.go

import (
	"app/base/models"
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
	testOrgID            = digestwriter.AccountNumber("1")
	workload             = digestwriter.Workload{}

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
	expectedSelect := `SELECT * FROM "image" WHERE (manifest_schema2_digest IN ($1) OR manifest_list_digest IN ($2) OR docker_image_digest IN ($3)) AND arch_id = $4`
	expectedSelect2 := `SELECT * FROM "cluster_image" WHERE cluster_id = $1`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`
	expectedCacheInsert := `UPDATE cluster SET cve_cache_critical = c.c, cve_cache_important = c.i, cve_cache_moderate = c.m, cve_cache_low = c.l FROM (SELECT COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $1 THEN cve.id ELSE NULL END), 0) AS c,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $2 THEN cve.id ELSE NULL END), 0) AS i,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $3 THEN cve.id ELSE NULL END), 0) AS m,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $4 THEN cve.id ELSE NULL END), 0) AS l FROM "cve" JOIN image_cve ON image_cve.cve_id = cve.id JOIN cluster_image ON cluster_image.image_id = image_cve.image_id WHERE cluster_image.cluster_id = $5) AS c WHERE cluster.id = $6`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest, firstDigest, firstDigest, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "manifest_schema2_digest", "manifest_list_digest", "docker_image_digest", "pyxis_id", "modified_date", "arch_id"}).
			AddRow(firstDigestID, firstDigest, firstDigest, firstDigest, 5, digestModifiedAtTime, 1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect2)).
		WithArgs(testClustedID).
		WillReturnRows(sqlmock.NewRows([]string{"cluster_id", "image_id"}))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(firstDigestID, 1))
	mock.ExpectExec(regexp.QuoteMeta(expectedCacheInsert)).
		WithArgs(models.Critical, models.Important, models.Moderate, models.Low, testClustedID, testClustedID).
		WillReturnResult(sqlmock.NewResult(testClustedID, 1))
	mock.ExpectCommit()

	// call the tested method
	err := digestwriter.LinkDigestsToCluster(storage, string(clusterName), testClustedID, 1, []string{firstDigest})
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
	expectedSelect := `SELECT * FROM "image" WHERE (manifest_schema2_digest IN ($1,$2) OR manifest_list_digest IN ($3,$4) OR docker_image_digest IN ($5,$6)) AND arch_id = $7`
	expectedSelect2 := `SELECT * FROM "cluster_image" WHERE cluster_id = $1`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`
	expectedCacheInsert := `UPDATE cluster SET cve_cache_critical = c.c, cve_cache_important = c.i, cve_cache_moderate = c.m, cve_cache_low = c.l FROM (SELECT COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $1 THEN cve.id ELSE NULL END), 0) AS c,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $2 THEN cve.id ELSE NULL END), 0) AS i,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $3 THEN cve.id ELSE NULL END), 0) AS m,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $4 THEN cve.id ELSE NULL END), 0) AS l FROM "cve" JOIN image_cve ON image_cve.cve_id = cve.id JOIN cluster_image ON cluster_image.image_id = image_cve.image_id WHERE cluster_image.cluster_id = $5) AS c WHERE cluster.id = $6`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest, secondDigest, firstDigest, secondDigest, firstDigest, secondDigest, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "manifest_schema2_digest", "manifest_list_digest", "docker_image_digest", "pyxis_id", "modified_date"}).
			AddRow(firstDigestID, firstDigest, firstDigest, firstDigest, 5, digestModifiedAtTime).
			AddRow(secondDigestID, secondDigest, secondDigest, secondDigest, 6, digestModifiedAtTime),
		)
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect2)).
		WithArgs(testClustedID).
		WillReturnRows(sqlmock.NewRows([]string{"cluster_id", "image_id"}))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(secondDigestID, 2))
	mock.ExpectExec(regexp.QuoteMeta(expectedCacheInsert)).
		WithArgs(models.Critical, models.Important, models.Moderate, models.Low, testClustedID, testClustedID).
		WillReturnResult(sqlmock.NewResult(testClustedID, 1))
	mock.ExpectCommit()

	// call the tested method
	err := digestwriter.LinkDigestsToCluster(storage, string(clusterName), testClustedID, 1, []string{firstDigest, secondDigest})
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
	err := storage.WriteClusterInfo(invalidClusterName, testOrgID, workload, []string{firstDigest})
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
	expectedSelect := `SELECT * FROM "account" WHERE "account"."org_id" = $1 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testAccountID))

	// Expect a SELECT FROM cluster statement since the previous select returned a record
	// Since it returns a record, there is no need to expect the INSERT INTO "cluster" statement
	expectedSelect = `SELECT "cluster"."id","cluster"."uuid","cluster"."account_id","cluster"."last_seen","cluster"."workload" FROM "cluster" WHERE "cluster"."uuid" = $1 AND "cluster"."account_id" = $2 AND "cluster"."last_seen" = $3 AND "cluster"."workload" = $4 ORDER BY "cluster"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(clusterName, testAccountID, anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "account_id"}).AddRow(testClustedID, clusterName, testAccountID))

	// Since it all went smoothly, the digest is linked to the cluster
	expectedSelect = `SELECT * FROM "arch" WHERE name = $1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs("amd64").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "amd64"))

	expectedSelect = `SELECT * FROM "image" WHERE (manifest_schema2_digest IN ($1) OR manifest_list_digest IN ($2) OR docker_image_digest IN ($3)) AND arch_id = $4`
	expectedSelect2 := `SELECT * FROM "cluster_image" WHERE cluster_id = $1`
	expectedInsert := `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`
	expectedCacheInsert := `UPDATE cluster SET cve_cache_critical = c.c, cve_cache_important = c.i, cve_cache_moderate = c.m, cve_cache_low = c.l FROM (SELECT COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $1 THEN cve.id ELSE NULL END), 0) AS c,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $2 THEN cve.id ELSE NULL END), 0) AS i,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $3 THEN cve.id ELSE NULL END), 0) AS m,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $4 THEN cve.id ELSE NULL END), 0) AS l FROM "cve" JOIN image_cve ON image_cve.cve_id = cve.id JOIN cluster_image ON cluster_image.image_id = image_cve.image_id WHERE cluster_image.cluster_id = $5) AS c WHERE cluster.id = $6`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest, firstDigest, firstDigest, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "manifest_schema2_digest", "manifest_list_digest", "docker_image_digest", "pyxis_id", "modified_date", "arch_id"}).
			AddRow(firstDigestID, firstDigest, firstDigest, firstDigest, 5, digestModifiedAtTime, 1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect2)).
		WithArgs(testClustedID).
		WillReturnRows(sqlmock.NewRows([]string{"cluster_id", "image_id"}))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WillReturnResult(sqlmock.NewResult(firstDigestID, 1))
	mock.ExpectExec(regexp.QuoteMeta(expectedCacheInsert)).
		WithArgs(models.Critical, models.Important, models.Moderate, models.Low, testClustedID, testClustedID).
		WillReturnResult(sqlmock.NewResult(testClustedID, 1))
	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteClusterInfo(clusterName, testOrgID, workload, []string{firstDigest})
	assert.Nil(t, err, "No error expected")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfoNoAccountForClusterName tests that
// Storage.WriteClusterInfo creates a new record in account table if
// no account can be retrieved for given clusterName and updates
// the cluster_image table properly
func TestWriteClusterInfoNoAccountForClusterName(t *testing.T) {
	storedAccountID := 10
	storedClusterID := 2

	storage, mock := NewMockStorage(t)
	mock.ExpectBegin()

	//Expect the select query and mock that it returns 0 records
	expectedSelect := `SELECT * FROM "account" WHERE "account"."org_id" = $1 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	//Expect an 'INSERT INTO account' statement since the select returned no rows
	expectedInsert := `INSERT INTO "account" ("org_id") VALUES ($1) ON CONFLICT DO NOTHING RETURNING "id"`
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsert)).WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storedAccountID))

	// Expect a SELECT FROM cluster statement
	// Mock that it would not return any row since there was no data for current account
	expectedSelect = `SELECT "cluster"."id","cluster"."uuid","cluster"."account_id","cluster"."last_seen","cluster"."workload" FROM "cluster" WHERE "cluster"."uuid" = $1 AND "cluster"."account_id" = $2 AND "cluster"."last_seen" = $3 AND "cluster"."workload" = $4 ORDER BY "cluster"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries, anyArgForMockSQLQueries, anyArgForMockSQLQueries, anyArgForMockSQLQueries).
		WillReturnRows(sqlmock.NewRows([]string{"id", "uuid", "account_id"}))
	//Expect an 'INSERT INTO cluster' with the ID of the created account record
	expectedInsert = `INSERT INTO "cluster" ("uuid","account_id","last_seen","workload") VALUES ($1,$2,$3,$4) ON CONFLICT ("uuid") DO UPDATE SET "uuid"="excluded"."uuid","account_id"="excluded"."account_id","last_seen"="excluded"."last_seen","workload"="excluded"."workload" RETURNING "id"`
	clusterUUID, _ := uuid.Parse(string(clusterName))
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsert)).WithArgs(clusterUUID, storedAccountID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(storedClusterID))

	// Since it all went smoothly, the digest is linked to the cluster
	expectedSelect = `SELECT * FROM "arch" WHERE name = $1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs("amd64").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "amd64"))

	expectedSelect = `SELECT * FROM "image" WHERE (manifest_schema2_digest IN ($1) OR manifest_list_digest IN ($2) OR docker_image_digest IN ($3)) AND arch_id = $4`
	expectedSelect2 := `SELECT * FROM "cluster_image" WHERE cluster_id = $1`
	expectedInsert = `INSERT INTO "cluster_image" ("cluster_id","image_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`
	expectedCacheInsert := `UPDATE cluster SET cve_cache_critical = c.c, cve_cache_important = c.i, cve_cache_moderate = c.m, cve_cache_low = c.l FROM (SELECT COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $1 THEN cve.id ELSE NULL END), 0) AS c,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $2 THEN cve.id ELSE NULL END), 0) AS i,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $3 THEN cve.id ELSE NULL END), 0) AS m,
				COALESCE(COUNT(DISTINCT CASE WHEN cve.severity = $4 THEN cve.id ELSE NULL END), 0) AS l FROM "cve" JOIN image_cve ON image_cve.cve_id = cve.id JOIN cluster_image ON cluster_image.image_id = image_cve.image_id WHERE cluster_image.cluster_id = $5) AS c WHERE cluster.id = $6`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(firstDigest, firstDigest, firstDigest, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "manifest_schema2_digest", "manifest_list_digest", "docker_image_digest", "pyxis_id", "modified_date", "arch_id"}).
			AddRow(firstDigestID, firstDigest, firstDigest, firstDigest, 5, digestModifiedAtTime, 1))

	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect2)).
		WithArgs(storedClusterID).
		WillReturnRows(sqlmock.NewRows([]string{"cluster_id", "image_id"}))

	mock.ExpectExec(regexp.QuoteMeta(expectedInsert)).WithArgs(storedClusterID, firstDigestID).
		WillReturnResult(sqlmock.NewResult(firstDigestID, 1))

	mock.ExpectExec(regexp.QuoteMeta(expectedCacheInsert)).
		WithArgs(models.Critical, models.Important, models.Moderate, models.Low, storedClusterID, storedClusterID).
		WillReturnResult(sqlmock.NewResult(testClustedID, 1))

	mock.ExpectCommit()

	// call the tested method
	err := storage.WriteClusterInfo(clusterName, testOrgID, workload, []string{firstDigest})
	assert.Nil(t, err, "No error expected")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

// TestWriteClusterInfoErrorWritingAccount tests that
// Storage.WriteClusterInfo aborts and DB is rolled back to
// previous state if an error happens while getting or creating
// the account info
func TestWriteClusterInfoErrorWritingAccount(t *testing.T) {
	storage, mock := NewMockStorage(t)
	mock.ExpectBegin()

	// expect select from account statement
	expectedSelect := `SELECT * FROM "account" WHERE "account"."org_id" = $1 ORDER BY "account"."id" LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelect)).
		WithArgs(anyArgForMockSQLQueries).
		WillReturnError(sql.ErrConnDone)

	mock.ExpectRollback()

	// call the tested method
	err := storage.WriteClusterInfo(clusterName, testOrgID, workload, []string{firstDigest})
	assert.Error(t, err, "account table shouldn't be updated and WriteClusterInfo should abort")

	// check that no SQL operations are done
	checkAllExpectations(t, mock)
}

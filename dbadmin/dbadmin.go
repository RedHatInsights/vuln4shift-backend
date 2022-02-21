package dbadmin

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" // We load migrations from local folder

	"app/base/utils"
)

var migrationFiles = "file://./dbadmin/migrations"

type logger struct{}

func (t logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (t logger) Verbose() bool {
	return true
}

func setDbUserPassword(conn *sql.DB, user string, passEnvVar string) {
	if password, ok := os.LookupEnv(passEnvVar); !ok {
		log.Fatalf("Unable to get env var: %s.\n", passEnvVar)
	} else {
		if _, err := conn.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", user, password)); err != nil {
			log.Printf("Setting password failed: %s", err) // Log but do not fail if user doesn't exist
		}
	}
}

func setDbUsersPasswords(conn *sql.DB) {
	setDbUserPassword(conn, "archive_db_writer", "USER_ARCHIVE_DB_WRITER_PASS")
	setDbUserPassword(conn, "pyxis_gatherer", "USER_PYXIS_GATHERER_PASS")
	setDbUserPassword(conn, "vmaas_gatherer", "USER_VMAAS_GATHERER_PASS")
	setDbUserPassword(conn, "cve_aggregator", "USER_CVE_AGGREGATOR_PASS")
	setDbUserPassword(conn, "manager", "USER_MANAGER_PASS")
}

func Start() {
	conn, err := utils.GetStandardDbConnection()
	if err != nil {
		log.Fatalf("Unable to connect to database: %s\n", err)
	}
	defer conn.Close()

	var driver database.Driver
	for { // Wait until DB is ready
		driver, err = pgx.WithInstance(conn, &pgx.Config{})
		if err != nil {
			log.Printf("Unable to get database driver, retrying: %s\n", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	m, err := migrate.NewWithDatabaseInstance(migrationFiles, utils.GetEnv("POSTGRES_DB", ""), driver)
	if err != nil {
		log.Fatalf("Unable to get migration interface: %s\n", err)
	}

	m.Log = logger{} // Set custom logger

	schemaMigration := utils.GetEnv("SCHEMA_MIGRATION", "-1") // Check env variable to migrate to specific version
	schemaMigrationInt, err := strconv.Atoi(schemaMigration)
	if err != nil {
		log.Fatalf("Unable to convert to int: %s\n", schemaMigration)
	}

	if schemaMigrationInt < 0 {
		err = m.Up() // Upgrade to the latest
	} else {
		err = m.Migrate(uint(schemaMigrationInt)) // Upgrade/Downgrade to the specific version
	}

	if err != nil {
		if err.Error() == "no change" {
			log.Printf("Schema is up to date.")
		} else {
			log.Fatalf("Error runnning the migration: %s", err.Error())
		}
	}

	setDbUsersPasswords(conn)
}

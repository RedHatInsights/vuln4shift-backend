package utils

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // Needed to make pgx work with database/sql
)

func GetDbURL() string {
	dbHost := GetEnv("POSTGRES_HOST", "vuln4shift_database")
	dbPort := GetEnv("POSTGRES_PORT", "5432")
	dbName := GetEnv("POSTGRES_DB", "vuln4shift")
	dbUser := GetEnv("POSTGRES_USER", "vuln4shift_admin")
	dbPassword := GetEnv("POSTGRES_PASSWORD", "vuln4shift_admin_pwd")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func GetDbConnection() (*pgx.Conn, error) {
	dbURL := GetDbURL()
	conn, err := pgx.Connect(context.Background(), dbURL)
	return conn, err
}

func GetStandardDbConnection() (*sql.DB, error) {
	dbURL := GetDbURL()
	conn, err := sql.Open("pgx", dbURL)
	return conn, err
}

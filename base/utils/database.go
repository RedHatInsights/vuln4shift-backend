package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // Needed to make pgx work with database/sql
)

func getDbVar(envVar string) string {
	value, ok := os.LookupEnv(envVar)
	if !ok {
		log.Fatalf("Unable to get env var: %s.\n", envVar)
	}
	return value
}

func GetDbURL() string {
	dbHost := getDbVar("POSTGRES_HOST")
	dbPort := getDbVar("POSTGRES_PORT")
	dbName := getDbVar("POSTGRES_DB")
	dbUser := getDbVar("POSTGRES_USER")
	dbPassword := getDbVar("POSTGRES_PASSWORD")
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

package utils

import (
	"app/base/models"
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"gorm.io/gorm"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // Needed to make pgx work with database/sql
)

func GetDbURL(admin bool) string {
	var dbUser, dbPassword string
	if admin {
		dbUser = Cfg.DbAdminUser
		dbPassword = Cfg.DbAdminPassword
	} else {
		dbUser = Cfg.DbUser
		dbPassword = Cfg.DbPassword
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbUser, url.QueryEscape(dbPassword), Cfg.DbHost, Cfg.DbPort, Cfg.DbName)
}

func GetDbConnection(admin bool) (*pgx.Conn, error) {
	dbURL := GetDbURL(admin)
	conn, err := pgx.Connect(context.Background(), dbURL)
	return conn, err
}

func GetStandardDbConnection(admin bool) (*sql.DB, error) {
	dbURL := GetDbURL(admin)
	conn, err := sql.Open("pgx", dbURL)
	return conn, err
}

func DbConfigure() (*gorm.DB, error) {
	dsn := GetDbURL(false)

	db, err := models.GetGormConnection(dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

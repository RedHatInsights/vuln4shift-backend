package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dbUser      string
	dbPswd      string
	dbAdminUser string
	dbAdminPswd string
	dbPort      int
	dbName      string
	dbHost      string
)

func setupEnvVars() {
	Cfg.DbUser = "test-user"
	Cfg.DbPassword = "test-password"
	Cfg.DbAdminUser = "test-admin-user"
	Cfg.DbAdminPassword = "test-admin-password"
	Cfg.DbPort = 25432
	Cfg.DbName = "vuln4shift"
	Cfg.DbHost = "0.0.0.0"
}

func rollbackEnvVars() {
	Cfg.DbUser = dbUser
	Cfg.DbPassword = dbPswd
	Cfg.DbAdminUser = dbAdminUser
	Cfg.DbAdminPassword = dbAdminPswd
	Cfg.DbPort = dbPort
	Cfg.DbName = dbName
	Cfg.DbHost = dbHost
}

func TestGetDbURL(t *testing.T) {
	setupEnvVars()
	defer rollbackEnvVars()

	expectedURL := "postgresql://test-user:test-password@0.0.0.0:25432/vuln4shift"
	assert.Equal(t, expectedURL, GetDbURL(false))
}

func TestGetDbURLAdmin(t *testing.T) {
	setupEnvVars()
	defer rollbackEnvVars()

	expectedURL := "postgresql://test-admin-user:test-admin-password@0.0.0.0:25432/vuln4shift"
	assert.Equal(t, expectedURL, GetDbURL(true))
}

func TestGetDbConnection(t *testing.T) {
	rollbackEnvVars()

	conn, err := GetDbConnection(false)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestGetStandardDbConnection(t *testing.T) {
	conn, err := GetStandardDbConnection(false)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestGetStandardDbConnectionAdmin(t *testing.T) {
	conn, err := GetStandardDbConnection(true)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestMain(m *testing.M) {
	dbUser = Cfg.DbUser
	dbPswd = Cfg.DbPassword
	dbAdminUser = Cfg.DbAdminUser
	dbAdminPswd = Cfg.DbAdminPassword
	dbPort = Cfg.DbPort
	dbName = Cfg.DbName
	dbHost = Cfg.DbHost
	os.Exit(m.Run())
}

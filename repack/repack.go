package repack

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"app/base/utils"
)

var (
	logger *logrus.Logger
	tables = map[string]string{
		"cluster":          "account_id,uuid",
		"cluster_image":    "cluster_id",
		"image":            "manifest_schema2_digest,manifest_list_digest,docker_image_digest",
		"image_cve":        "cve_id",
		"repository":       "registry,repository",
		"repository_image": "repository_id",
	}
	pgRepackArgs = []string{
		"--no-superuser-check",
		"--no-password",
		"-d", utils.Cfg.DbName,
		"-h", utils.Cfg.DbHost,
		"-p", fmt.Sprintf("%d", utils.Cfg.DbPort),
		"-U", utils.Cfg.DbAdminUser,
	}
)

func init() {
	var err error
	logger, err = utils.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// GetCmd returns command that calls pg_repack with table.
// Args are appended to the necessary pgRepackArgs. Stdout and stderr of the subprocess are redirected to host.
func getCmd(table string, args ...string) *exec.Cmd {
	fullArgs := pgRepackArgs
	fullArgs = append(fullArgs, "-I", table)
	fullArgs = append(fullArgs, args...)
	cmd := exec.Command("pg_repack", fullArgs...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", utils.Cfg.DbAdminPassword))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Repack runs pg_repack with table. If columns are provided, cluster by these columns is executed as well.
func repack(table string, columns string) error {
	var clusterCmd *exec.Cmd
	if len(columns) > 0 {
		clusterCmd = getCmd(table, "-o", columns)
	} else {
		clusterCmd = getCmd(table)
	}
	err := clusterCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func createPgRepackExtension() error {
	logger.Info("Checking pg_repack extension in DB")
	conn, err := utils.GetStandardDbConnection(true)
	if err != nil {
		return err
	}
	defer conn.Close()

	var defaultVersion, installedVersion string
	query := "SELECT COALESCE(default_version, ''), COALESCE(installed_version, '') FROM pg_available_extensions WHERE name = 'pg_repack'"
	if err := conn.QueryRow(query).Scan(&defaultVersion, &installedVersion); err != nil {
		logger.Error("Unable to get available pg_repack extension version, it may not be available")
		return err
	}

	if installedVersion != defaultVersion {
		if installedVersion != "" {
			logger.Infof("Dropping existing pg_repack extension version %s", installedVersion)
			if _, err := conn.Exec("DROP EXTENSION pg_repack"); err != nil {
				logger.Error("Dropping extension failed")
				return err
			}
		}
		logger.Infof("Creating pg_repack extension version %s", defaultVersion)
		if _, err := conn.Exec("CREATE EXTENSION pg_repack"); err != nil {
			logger.Error("Creating extension failed")
			return err
		}
	} else {
		logger.Infof("pg_repack extension version %s already exists", installedVersion)
	}
	return nil
}

func Start() {
	logger.Info("Starting pg_repack job.")
	err := createPgRepackExtension()
	if err != nil {
		logger.Fatalf("Failed to create pg_repack extension, error: %s", err.Error())
	}
	for table, columns := range tables {
		err := repack(table, columns)
		if err != nil {
			logger.Fatalf("Failed to repack table %s, error: %s", table, err.Error())
		}
		logger.Infof("Successfully repacked table %s", table)
	}
}

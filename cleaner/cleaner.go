package cleaner

import (
	"app/base/logging"
	"app/base/utils"
	"fmt"
	"os"
)

// Cleaner interface
type Cleaner interface {
	Clean() error
}

// Start starts the cleaner job
func Start() {
	logger, err := logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger")
		os.Exit(1)
	}

	clusterRetention := utils.GetEnv("CLUSTER_RETENTION_DAYS", 0)
	if !(clusterRetention > 0) {
		logger.Fatalf("CLUSTER_RETENTION_DAYS env not set")
	}

	clusterCleaner, err := NewClusterCleaner(uint(clusterRetention))
	if err != nil {
		logger.Fatalf("Error creating cluster cleaner: %s", err)
	}

	err = clusterCleaner.Clean()
	if err != nil {
		logger.Errorf("Error deleting clusters: %s", err)
	}
}

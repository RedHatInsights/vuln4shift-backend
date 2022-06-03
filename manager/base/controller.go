package base

import (
	"app/base/logging"
	"app/base/utils"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Controller struct {
	Conn   *gorm.DB
	Logger *logrus.Logger
}

func CreateControllerLogger() *logrus.Logger {
	logger, err := logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}

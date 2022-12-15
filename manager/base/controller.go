package base

import (
	"app/base/utils"
	"app/manager/amsclient"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Controller struct {
	Conn      *gorm.DB
	AMSClient amsclient.AMSClient
	Logger    *logrus.Logger
}

func CreateControllerLogger() *logrus.Logger {
	logger, err := utils.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}

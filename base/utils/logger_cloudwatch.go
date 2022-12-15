package utils

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	lc "github.com/redhatinsights/platform-go-middlewares/logging/cloudwatch"
	"github.com/sirupsen/logrus"
)

var hook *lc.Hook

// Try to init CloudWatch logging
func trySetupCloudWatchLogging(logger *logrus.Logger) {
	if hook == nil {
		if Cfg.CloudWatchAccessKeyID == "" {
			logger.Info("config for aws CloudWatch not loaded")
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			logger.Errorf("unable to get hostname to set CloudWatch log_stream: %s", err.Error())
			return
		}

		cred := credentials.NewStaticCredentials(Cfg.CloudWatchAccessKeyID, Cfg.CloudWatchSecretAccesskey, "")
		awsconf := aws.NewConfig().WithRegion(Cfg.CloudWatchRegion).WithCredentials(cred)
		hook, err = lc.NewBatchingHook(Cfg.CloudWatchLogGroup, hostname, awsconf, 10*time.Second)
		if err != nil {
			logger.Errorf("unable to setup CloudWatch logging: %s", err.Error())
			return
		}
		logger.Info("CloudWatch logging configured")
	}
	logger.AddHook(hook)
}

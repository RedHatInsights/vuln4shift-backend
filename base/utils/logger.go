package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// ParseLoglevel parses string logging level value to Walrus logging values.
func ParseLoglevel(level string) (logrus.Level, error) {
	level = strings.ToUpper(level)
	switch level {
	case "TRACE":
		return logrus.TraceLevel, nil
	case "DEBUG":
		return logrus.DebugLevel, nil
	case "INFO":
		return logrus.InfoLevel, nil
	case "WARN":
		return logrus.WarnLevel, nil
	case "ERROR":
		return logrus.ErrorLevel, nil
	case "PANIC":
		return logrus.PanicLevel, nil
	case "FATAL":
		return logrus.FatalLevel, nil
	default:
		return logrus.InfoLevel, errors.New("invalid loglevel given")
	}
}

// CreateLogger creates walrus type logger based on logging level.
func CreateLogger(level string) (*logrus.Logger, error) {
	logger := logrus.New()

	logrusLevel, err := ParseLoglevel(level)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(logrusLevel)
	logger.SetOutput(os.Stdout)
	trySetupCloudWatchLogging(logger)
	return logger, nil
}

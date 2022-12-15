package digestwriter

// Entry point of the digestwriter package

import (
	"app/base/utils"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota
	// ExitStatusStorageError is returned in case of any consumer-related error
	ExitStatusStorageError
	// ExitStatusConsumerError is returned in case of any consumer-related error
	ExitStatusConsumerError
)

func setupLogger() {
	if logger == nil {
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
}

// startConsumer function starts the Kafka consumer.
func startConsumer(storage Storage) (*utils.KafkaConsumer, error) {
	consumer, err := NewConsumer(storage)
	if err != nil {
		return nil, err
	}
	consumer.Serve()
	return consumer, nil
}

// Start function tries to start the digest writer service.
func Start() {
	setupLogger()
	logger.Infoln("Initializing digest writer...")

	RunMetrics()

	storage, err := NewStorage()
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing storage")
		logger.Logln(logrus.FatalLevel, err.Error())
		logger.Exit(ExitStatusStorageError)
	}

	consumer, err := startConsumer(storage)
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing consumer")
		logger.Logln(logrus.FatalLevel, err.Error())
		logger.Exit(ExitStatusConsumerError)
	}
	defer consumer.Close()

	logger.Infoln("Digest writer done")
	logger.Exit(ExitStatusOK)
}

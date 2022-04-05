package digestwriter

// Entry point of the digestwriter package

import (
	"app/base/logging"
	"app/base/utils"

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
	var err error
	logger, err = logging.CreateLogger(utils.GetEnv("LOGGING_LEVEL", "INFO"))
	if err != nil {
		panic("couldn't set up logger with given LOGGING_LEVEL environment variable nor default (INFO)")
	}
}

// startConsumer function starts the Kafka consumer.
func startConsumer(storage Storage) (*KafkaConsumer, error) {
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

	storage, err := NewStorage()
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing storage")
		logger.Exit(ExitStatusStorageError)
	}

	consumer, err := startConsumer(storage)
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing consumer")
		logger.Exit(ExitStatusConsumerError)
	}
	defer consumer.Close()

	logger.Infoln("Digest writer done")
	logger.Exit(ExitStatusOK)
}

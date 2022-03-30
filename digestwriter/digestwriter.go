package digestwriter

import (
	"app/base/logging"
	"app/base/utils"

	"github.com/sirupsen/logrus"
)

const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota
	// ExitStatusStorageError is returned in case of any consumer-related error
	ExitStatusStorageError
	// ExitStatusConsumerError is returned in case of any consumer-related error
	ExitStatusConsumerError
)

// startConsumer function starts the Kafka consumer.
func startConsumer(storage Storage, logger *logrus.Logger) (*KafkaConsumer, error) {
	consumer, err := NewConsumer(storage, logger)
	if err != nil {
		return nil, err
	}
	consumer.Serve()
	return consumer, nil
}

// Start function tries to start the digest writer service.
func Start() {
	logger, err := logging.CreateLogger(utils.Getenv("LOGGING_LEVEL", "INFO"))
	if err != nil {
		panic("Invalid LOGGING_LEVEL environment variable set")
	}
	logger.Infoln("Initializing digest writer...")

	storage, err := NewStorage(logger)
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing storage")
		logger.Exit(ExitStatusStorageError)
	}

	consumer, err := startConsumer(storage, logger)
	if err != nil {
		logger.Logln(logrus.FatalLevel, "Error initializing consumer")
		logger.Exit(ExitStatusConsumerError)
	}
	defer consumer.Close()

	logger.Infoln("Digest writer done")
	logger.Exit(ExitStatusOK)
}

package utils

import (
	"app/base/logging"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

const (
	// key for topic name used in structured log messages
	topicKey = "topic"
	// key for message offset used in structured log messages
	offsetKey = "offset"
	// key for message partition used in structured log messages
	partitionKey = "partition"
	// key for error message used in structured log messages
	errorKey = "error"
	// key for duration of message processing used in structured log messages
	processingDurationKey = "processing_duration"
	// key for new message received timestamp
	messageTimestamp = "message_timestamp"
)

var (
	logger *logrus.Logger
)

func setupLogger() {
	if logger == nil {
		var err error
		logger, err = logging.CreateLogger(Cfg.LoggingLevel)
		if err != nil {
			fmt.Println("Error setting up logger.")
			os.Exit(1)
		}
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// Consumer interface for a topic consumer for any message broker
type Consumer interface {
	Serve()
	Close() error
	handleMessage(msg *sarama.ConsumerMessage)
}

// KafkaConsumerConfig configuration for connecting with Kafka broker
type KafkaConsumerConfig struct {
	// Address broker's address in <host>:<port> format
	Address string
	// IncomingTopic name of Kafka topic to consume from
	IncomingTopic string
	// Group name of the Kafka consumer group, if any
	Group string
}

// Processor interface that should be fulfilled by the struct that want to use
// the KafkaConsumer as "base" struct
type Processor interface {
	ProcessMessage(msg *sarama.ConsumerMessage) error
}

// KafkaConsumer holds the necessary properties for connecting to a Kafka topic,
// consuming and processing the messages, as well as basic metrics.
// It implements the Consumer interface defined in this source file, as well as
// the ConsumerGroupHandler interface defined in the sarama package
type KafkaConsumer struct {
	Config                               *KafkaConsumerConfig
	ConsumerGroup                        sarama.ConsumerGroup
	numberOfSuccessfullyConsumedMessages uint64
	numberOfErrorsConsumingMessages      uint64
	Ready                                chan bool
	Cancel                               context.CancelFunc
	Processor                            Processor
}

// NewKafkaConsumer constructs new implementation of KafkaConsumer, using
// the default sarama config if none is provided
func NewKafkaConsumer(saramaConfig *sarama.Config, processor Processor) (*KafkaConsumer, error) {
	setupLogger()
	if Cfg.KafkaBrokerAddress == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_ADDRESS")
	}
	if Cfg.KafkaBrokerConsumerGroup == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_CONSUMER_GROUP")
	}
	if Cfg.KafkaBrokerIncomingTopic == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_INCOMING_TOPIC")
	}
	if saramaConfig == nil {
		saramaConfig = sarama.NewConfig()
		saramaConfig.Version = sarama.V0_10_2_0

		timeout, err := time.ParseDuration(Cfg.KafkaConsumerTimeout)
		if err == nil && timeout != 0 {
			saramaConfig.Net.DialTimeout = timeout
			saramaConfig.Net.ReadTimeout = timeout
			saramaConfig.Net.WriteTimeout = timeout
		}
	}

	consumerGroup, err := sarama.NewConsumerGroup([]string{Cfg.KafkaBrokerAddress}, Cfg.KafkaBrokerConsumerGroup, saramaConfig)
	if err != nil {
		return nil, err
	}

	consumer := &KafkaConsumer{
		Config: &KafkaConsumerConfig{
			Address:       Cfg.KafkaBrokerAddress,
			IncomingTopic: Cfg.KafkaBrokerIncomingTopic,
			Group:         Cfg.KafkaBrokerConsumerGroup,
		},
		ConsumerGroup:                        consumerGroup,
		numberOfSuccessfullyConsumedMessages: 0,
		numberOfErrorsConsumingMessages:      0,
		Ready:                                make(chan bool),
		Processor:                            processor,
	}

	return consumer, nil
}

// Serve starts listening for messages and processing them.
// Must be called to consume the messages
func (consumer *KafkaConsumer) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	consumer.Cancel = cancel

	go func() {
		for {
			if err := consumer.ConsumerGroup.Consume(
				ctx, []string{consumer.Config.IncomingTopic}, consumer); err != nil {
				logger.WithFields(logrus.Fields{
					errorKey: err.Error(),
				}).Errorln("Unable to recreate Kafka session")
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logger.WithFields(logrus.Fields{
					errorKey: ctx.Err(),
				}).Infoln("stopping consumer")
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	// Wait until the consumer is ready
	logger.Debugln("waiting for kafka consumer to become ready")
	<-consumer.Ready
	logger.Debugln("kafka consumer is ready")

	// Actual processing is done in goroutine created by sarama (see ConsumeClaim below)
	logger.Infoln("started serving consumer")
	<-ctx.Done()
	logger.Infoln("sarama context cancelled, exiting")

	cancel()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.Infoln("kafka session has been terminated")
	return nil
}

// ConsumeClaim implements the message consuming loop
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.WithFields(logrus.Fields{
		offsetKey: claim.InitialOffset(),
	}).Debugln("Starting messages loop")

	for message := range claim.Messages() {
		consumer.handleMessage(message)
		session.MarkMessage(message, "")
	}

	return nil
}

// Close method closes all resources used by consumer
func (consumer *KafkaConsumer) Close() error {
	if consumer.Cancel != nil {
		consumer.Cancel()
	}

	if consumer.ConsumerGroup != nil {
		if err := consumer.ConsumerGroup.Close(); err != nil {
			logger.WithFields(logrus.Fields{
				errorKey: err.Error(),
			}).Errorln("Unable to close consumer group")
		}
	}
	return nil
}

// handleMessage handles the message and does basic logging and metrics update.
// It calls ProcessMessage, which must be implemented by the Processor defined
func (consumer *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage) {
	if consumer.Processor == nil {
		panic("message processor has not been set up. Aborting handleMessage")
	}
	if msg == nil {
		logger.Debugln("nil message")
		return
	}

	logger.WithFields(logrus.Fields{
		offsetKey:        msg.Offset,
		partitionKey:     msg.Partition,
		topicKey:         msg.Topic,
		messageTimestamp: msg.Timestamp,
	}).Debugln("Start processing incoming message")

	startTime := time.Now()
	err := consumer.Processor.ProcessMessage(msg)
	timeAfterProcessingMessage := time.Now()

	// Something went wrong while processing the message.
	if err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err.Error(),
		}).Errorln("Error processing the message consumed from Kafka")
		consumer.IncrementNumberOfErrorsConsumingMessages()
		/* ConsumingErrors.Inc() */
		return
	}

	logger.WithFields(logrus.Fields{
		offsetKey:             msg.Offset,
		partitionKey:          msg.Partition,
		topicKey:              msg.Topic,
		processingDurationKey: timeAfterProcessingMessage.Sub(startTime).Seconds(),
	}).Debugln("Processed incoming message successfully")
	consumer.IncrementNumberOfSuccessfullyConsumedMessages()
	/*ConsumedMessages.Inc()*/
}

// GetNumberOfSuccessfullyConsumedMessages returns number of consumed messages
// since creating KafkaConsumer obj
func (consumer *KafkaConsumer) GetNumberOfSuccessfullyConsumedMessages() uint64 {
	return consumer.numberOfSuccessfullyConsumedMessages
}

// IncrementNumberOfSuccessfullyConsumedMessages increments number of consumed messages
func (consumer *KafkaConsumer) IncrementNumberOfSuccessfullyConsumedMessages() {
	consumer.numberOfSuccessfullyConsumedMessages++
}

// GetNumberOfErrorsConsumingMessages returns number of errors during consuming messages
// since creating KafkaConsumer obj
func (consumer *KafkaConsumer) GetNumberOfErrorsConsumingMessages() uint64 {
	return consumer.numberOfErrorsConsumingMessages
}

// IncrementNumberOfErrorsConsumingMessages increments number of errors during consuming messages
func (consumer *KafkaConsumer) IncrementNumberOfErrorsConsumingMessages() {
	consumer.numberOfErrorsConsumingMessages++
}

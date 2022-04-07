package digestwriter

import (
	"app/base/utils"
	"context"
	"encoding/json"
	"errors"
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
)

// ClusterName represents the external cluster UUID contained in the consumed message
type ClusterName string

// JSONContent represents the workload info contained in the consumed message
type JSONContent map[string]*json.RawMessage

// IncomingMessage data structure is representation of message consumed from
// the configured topic
type IncomingMessage struct {
	Pods       int          `json:"-"`
	ImageCount int          `json:"imageCount"`
	Digests    *JSONContent `json:"images"`
	Namespaces *JSONContent `json:"-"`
}

// Consumer interface for a topic consumer for any message broker
type Consumer interface {
	Serve()
	Close() error
	ProcessMessage(msg *sarama.ConsumerMessage) error
}

// KafkaConsumerConfig represents the configuration for communicating
// with Kafka broker
type KafkaConsumerConfig struct {
	// Address represents Kafka address
	Address string
	// IncomingTopic is name of Kafka topic to consume from
	IncomingTopic string
	// Group is name of Kafka consumer group
	Group string
}

type KafkaConsumer struct {
	Config                               KafkaConsumerConfig
	ConsumerGroup                        sarama.ConsumerGroup
	Storage                              Storage
	numberOfSuccessfullyConsumedMessages uint64
	numberOfErrorsConsumingMessages      uint64
	Ready                                chan bool
	Cancel                               context.CancelFunc
}

// DefaultSaramaConfig is a config which will be used by default
// here you can use specific version of a protocol for example
// useful for testing
var DefaultSaramaConfig *sarama.Config

// NewConsumer constructs new implementation of Consumer interface
func NewConsumer(storage Storage) (*KafkaConsumer, error) {
	return NewWithSaramaConfig(DefaultSaramaConfig, storage)
}

// NewWithSaramaConfig constructs new implementation of Consumer interface with custom sarama config
func NewWithSaramaConfig(
	saramaConfig *sarama.Config,
	storage Storage,
) (*KafkaConsumer, error) {

	brokerAddress := utils.Getenv("KAFKA_BROKER_ADDRESS", "")
	if brokerAddress == "" {
		logger.Errorln("Unable to get env var: KAFKA_BROKER_ADDRESS")
	}
	group := utils.Getenv("KAFKA_BROKER_CONSUMER_GROUP", "")
	if group == "" {
		logger.Errorln("Unable to get env var: KAFKA_BROKER_CONSUMER_GROUP")
	}
	topic := utils.Getenv("KAFKA_BROKER_INCOMING_TOPIC", "")
	if topic == "" {
		logger.Errorln("Unable to get env var: KAFKA_BROKER_INCOMING_TOPIC")
	}
	if saramaConfig == nil {
		saramaConfig = sarama.NewConfig()
		saramaConfig.Version = sarama.V0_10_2_0

		timeout, err := time.ParseDuration(utils.Getenv("KAFKA_CONSUMER_TIMEOUT", ""))
		if err == nil && timeout != 0 {
			saramaConfig.Net.DialTimeout = timeout
			saramaConfig.Net.ReadTimeout = timeout
			saramaConfig.Net.WriteTimeout = timeout
		}
	}

	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerAddress}, group, saramaConfig)
	if err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err,
		}).Errorln("Couldn't setup Kafka consumer group with given config")
	}

	consumer := &KafkaConsumer{
		Config: KafkaConsumerConfig{
			Address:       brokerAddress,
			IncomingTopic: topic,
			Group:         group,
		},
		ConsumerGroup:                        consumerGroup,
		Storage:                              storage,
		numberOfSuccessfullyConsumedMessages: 0,
		numberOfErrorsConsumingMessages:      0,
		Ready:                                make(chan bool),
	}

	return consumer, nil
}

// Serve starts listening for messages and processing them. It blocks current thread.
func (consumer *KafkaConsumer) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	consumer.Cancel = cancel

	go func() {
		for {
			if err := consumer.ConsumerGroup.Consume(ctx, []string{consumer.Config.IncomingTopic}, consumer); err != nil {
				logger.WithFields(logrus.Fields{
					errorKey: err,
				}).Errorln("Unable to recreate Kafka session")
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logger.WithFields(logrus.Fields{
					errorKey: ctx.Err(),
				}).Errorln("Stopping consumer")
				return
			}

			logger.Info("Created new kafka session")

			consumer.Ready = make(chan bool)
		}
	}()

	// Wait until the consumer is ready
	logger.Info("Waiting for consumer to become ready")
	<-consumer.Ready
	logger.Info("Consumer is ready")

	// Actual processing is done in goroutine created by sarama (see ConsumeClaim below)
	logger.Info("Started serving consumer")
	<-ctx.Done()
	logger.Info("Context cancelled, exiting")

	cancel()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	logger.Info("New session has been setup")
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.Info("New session has been terminated")
	return nil
}

// ConsumeClaim starts a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.WithFields(logrus.Fields{
		offsetKey: claim.InitialOffset(),
	}).Info("Starting messages loop")

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
				errorKey: err,
			}).Errorln("Unable to close consumer group")
		}
	}

	return nil
}

// GetNumberOfSuccessfullyConsumedMessages returns number of consumed messages
// since creating KafkaConsumer obj
func (consumer *KafkaConsumer) GetNumberOfSuccessfullyConsumedMessages() uint64 {
	return consumer.numberOfSuccessfullyConsumedMessages
}

// GetNumberOfErrorsConsumingMessages returns number of errors during consuming messages
// since creating KafkaConsumer obj
func (consumer *KafkaConsumer) GetNumberOfErrorsConsumingMessages() uint64 {
	return consumer.numberOfErrorsConsumingMessages
}

// handleMessage handles the message and does all logging, metrics, etc
func (consumer *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage) {
	if msg == nil {
		logger.Infoln("nil message")
		return
	}

	logger.WithFields(logrus.Fields{
		offsetKey:           msg.Offset,
		partitionKey:        msg.Partition,
		topicKey:            msg.Topic,
		"message_timestamp": msg.Timestamp,
	}).Infoln("Start processing incoming message")

	startTime := time.Now()
	err := consumer.ProcessMessage(msg)
	timeAfterProcessingMessage := time.Now()

	// Something went wrong while processing the message.
	if err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err,
		}).Errorln("Error processing the message consumed from Kafka")
		consumer.numberOfErrorsConsumingMessages++
		/* ConsumingErrors.Inc() */
		return
	}

	logger.WithFields(logrus.Fields{
		offsetKey:             msg.Offset,
		partitionKey:          msg.Partition,
		topicKey:              msg.Topic,
		processingDurationKey: timeAfterProcessingMessage.Sub(startTime).Seconds(),
	}).Infoln("Processed incoming message successfully")
	consumer.numberOfSuccessfullyConsumedMessages++
	/*ConsumedMessages.Inc()*/
}

// ProcessMessage processes an incoming message
func (consumer *KafkaConsumer) ProcessMessage(msg *sarama.ConsumerMessage) error {
	tStart := time.Now()

	// Step #1: parse the incoming message
	message, err := parseMessage(msg.Value)
	if err != nil {
		/* ParseIncomingMessageError.Inc() */
		return err
	}

	/* ParsedIncomingMessage.Inc() */

	logger.Infoln("Parsed incoming message successfully")

	if message.Digests == nil {
		logger.Infoln("No digests were retrieved from incoming message")
		return nil
	}

	// Step #2: get digests into a slice of strings
	digests := extractDigestsFromMessage(message.Digests)

	logger.Infof("Extracted digests: %d\n", len(digests))

	if message.ImageCount != len(digests) {
		logger.Warnf("Expected number of digests: %d; Extracted digests: %d\n", message.ImageCount, len(digests))
	}

	// Step #3: write the digests into storage
	err = consumer.Storage.WriteDigests(digests)
	if err != nil {
		logger.WithFields(logrus.Fields{
			errorKey: err,
		}).Errorln("Error writing digests to database")
		/* StoredMessagesError.Inc() */
		return err
	}

	/* StoredMessagesOk.Inc() */

	// Step #5: print durations of all previous steps
	logger.WithFields(logrus.Fields{
		processingDurationKey: time.Now().Sub(tStart).Seconds(),
	}).Infoln("Stored digests successfully")

	return nil
}

func extractDigestsFromMessage(content *JSONContent) (digests []string) {
	// get the digest of each item
	digests = make([]string, len(*content))
	index := 0
	// TBD: We lose the ordering from the JSON by looping this way. Check if it matters
	for digest := range *content {
		digests[index] = digest
		index++
	}
	return
}

// parseMessage tries to parse incoming message and verify all required attributes
func parseMessage(messageValue []byte) (IncomingMessage, error) {
	var deserialized IncomingMessage

	err := json.Unmarshal(messageValue, &deserialized)
	if err != nil {
		return deserialized, err
	}

	if deserialized.Digests == nil {
		return deserialized, errors.New("missing required attribute 'images'")
	}

	if len(*deserialized.Digests) == 0 {
		return deserialized, errors.New("received message does not contain any digest")
	}

	return deserialized, nil
}

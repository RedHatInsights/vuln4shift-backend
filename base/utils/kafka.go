package utils

import (
	"context"
	"time"

	"github.com/pkg/errors"

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

func SetupLogger() {
	if logger == nil {
		var err error
		logger, err = CreateLogger(Cfg.LoggingLevel)
		if err != nil {
			logFatalf("Error setting up logger: %s", err.Error())
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
	// Shutdown function called during Close operation
	Shutdown func()
}

// NewKafkaConsumer constructs new implementation of KafkaConsumer, using
// the default sarama config if none is provided
func NewKafkaConsumer(saramaConfig *sarama.Config, processor Processor) (*KafkaConsumer, error) {
	SetupLogger()
	if Cfg.KafkaBrokerAddress == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_ADDRESS")
	}
	if Cfg.KafkaBrokerConsumerGroup == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_CONSUMER_GROUP")
	}
	if Cfg.KafkaBrokerIncomingTopic == "" {
		return nil, errors.New("unable to get env var: KAFKA_BROKER_INCOMING_TOPIC")
	}
	if Cfg.KafkaPayloadTrackerTopic == "" {
		return nil, errors.New("unable to get env var: KAFKA_PAYLOAD_TRACKER_TOPIC")
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

	if Cfg.KafkaBroker.Sasl != nil {
		err := SetKafkaSSLConfig(saramaConfig)
		if err != nil {
			return nil, err
		}
	}

	if Cfg.KafkaBroker.Authtype != nil {
		err := SetKafkaTLSConfig(saramaConfig)
		if err != nil {
			return nil, err
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
	// Call optional shutdown operation if defined
	if sd := consumer.Shutdown; sd != nil {
		sd()
	}

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
		ConsumingErrors.Inc()
		return
	}

	logger.WithFields(logrus.Fields{
		offsetKey:             msg.Offset,
		partitionKey:          msg.Partition,
		topicKey:              msg.Topic,
		processingDurationKey: timeAfterProcessingMessage.Sub(startTime).Seconds(),
	}).Debugln("Processed incoming message successfully")
	consumer.IncrementNumberOfSuccessfullyConsumedMessages()
	ConsumedMessages.Inc()
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

// Writer interface for writing messages to Kafka topic
type Writer interface {
	// Input returns channel receiving messages to be sent to the topic
	Input() chan<- *sarama.ProducerMessage
	// Successes returns channel containing messages successfully sent
	Successes() <-chan *sarama.ProducerMessage
	// Errors returns channel containing write errors
	Errors() <-chan *sarama.ProducerError
	Close() error
}

// Producer interface for publishing Kafka messages
type Producer interface {
	SendMessage(key, value sarama.Encoder)
	Close()
}

// KafkaProducerConfig configuration for connecting with Kafka broker
type KafkaProducerConfig struct {
	// Address broker's address in <host>:<port> format
	Address string
	// Topic name of Kafka topic to consume from
	Topic string
}

type KafkaProducer struct {
	Config                               *KafkaProducerConfig
	numberOfSuccessfullyProducedMessages uint64
	numberOfErrorsProducingMessages      uint64
	Writer                               Writer
	Enqueued                             int
}

func NewKafkaProducer(saramaConfig *sarama.Config, address, topic string) (*KafkaProducer, error) {
	SetupLogger()

	if address == "" {
		return nil, errors.New("empty broker address")
	}
	if topic == "" {
		return nil, errors.New("empty producer topic")
	}
	if saramaConfig == nil {
		saramaConfig = sarama.NewConfig()
		saramaConfig.Version = sarama.V0_10_2_0

		timeout, err := time.ParseDuration(Cfg.KafkaProducerTimeout)
		if err == nil && timeout != 0 {
			saramaConfig.Net.DialTimeout = timeout
			saramaConfig.Net.ReadTimeout = timeout
			saramaConfig.Net.WriteTimeout = timeout
		}
	}

	if Cfg.KafkaBroker.Sasl != nil {
		err := SetKafkaSSLConfig(saramaConfig)
		if err != nil {
			return nil, err
		}
	}

	if Cfg.KafkaBroker.Authtype != nil {
		err := SetKafkaTLSConfig(saramaConfig)
		if err != nil {
			return nil, err
		}
	}

	writer, err := sarama.NewAsyncProducer([]string{address}, saramaConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new Sarama async producer")
	}

	producer := &KafkaProducer{
		Config: &KafkaProducerConfig{
			Address: address,
			Topic:   topic,
		},
		numberOfSuccessfullyProducedMessages: 0,
		numberOfErrorsProducingMessages:      0,
		Writer:                               writer,
		Enqueued:                             0,
	}

	return producer, nil
}

func (producer *KafkaProducer) awaitWriteResult() error {
	select {
	case msg := <-producer.Writer.Successes():
		logger.Debugf("successfully sent a message with a key=%s to the %s topic", msg.Key, producer.Config.Topic)
		producer.IncrementNumberOfSuccessfullyProducedMessages()
		producer.Enqueued--
		return nil
	case err := <-producer.Writer.Errors():
		producer.IncrementNumberOfErrorsProducingMessages()
		producer.Enqueued--
		return err
	}
}

func (producer *KafkaProducer) write(msg *sarama.ProducerMessage) error {
	select {
	case producer.Writer.Input() <- msg:
		logger.Debugf("enqueued new message with a key=%s on the %s topic", msg.Key, producer.Config.Topic)
		producer.Enqueued++
		return producer.awaitWriteResult()
	case err := <-producer.Writer.Errors():
		producer.IncrementNumberOfErrorsProducingMessages()
		producer.Enqueued--
		return err
	}
}

// SendMessage composes Sarama message and sends it to the topic waiting for returning success or error value
func (producer *KafkaProducer) SendMessage(key, value sarama.Encoder) {
	msg := &sarama.ProducerMessage{
		Topic:     producer.Config.Topic,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}

	logger.Debugf("attempting to send a new message with a key=%s on the %s topic", key, producer.Config.Topic)
	if err := producer.write(msg); err != nil {
		logger.Errorf("failed to write kafka message: %s", err.Error())
	}
}

// Close closes underlying Writer which must be called to not leak memory, logging errors for potentially lost enqueued messages in the process
func (producer *KafkaProducer) Close() {
	if producer.Enqueued > 0 {
		logger.Warnf("closing underlying Kafka writer with %d unprocessed messages", producer.Enqueued)
	}

	if err := producer.Writer.Close(); err != nil {
		logger.Errorf("errors occurred during closing underlying Kafka writer: %s", err.Error())
	}

	logger.Info("successfully closed kafka producer")
}

// GetNumberOfSuccessfullyProducedMessages returns number of produced messages
func (producer *KafkaProducer) GetNumberOfSuccessfullyProducedMessages() uint64 {
	return producer.numberOfSuccessfullyProducedMessages
}

// IncrementNumberOfSuccessfullyProducedMessages increments number of produced messages
func (producer *KafkaProducer) IncrementNumberOfSuccessfullyProducedMessages() {
	producer.numberOfSuccessfullyProducedMessages++
}

// GetNumberOfErrorsProducingMessages returns number of errors during producing messages
func (producer *KafkaProducer) GetNumberOfErrorsProducingMessages() uint64 {
	return producer.numberOfErrorsProducingMessages
}

// IncrementNumberOfErrorsProducingMessages increments number of errors during producing messages
func (producer *KafkaProducer) IncrementNumberOfErrorsProducingMessages() {
	producer.numberOfErrorsProducingMessages++
}

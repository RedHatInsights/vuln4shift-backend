package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

type TestProcessor struct {
	numberOfMessagesProcessed uint64
}

func (d *TestProcessor) GetNumberOfMessagesProcessed() uint64 {
	return d.numberOfMessagesProcessed
}
func (d *TestProcessor) IncrementNumberOfMessagesProcessed() {
	d.numberOfMessagesProcessed++
}

// ProcessMessage processes an incoming message if provided or returns an error
func (d *TestProcessor) ProcessMessage(msg *sarama.ConsumerMessage) error {
	if msg == nil {
		return errors.New("couldn't process message")
	}
	d.IncrementNumberOfMessagesProcessed()
	return nil
}

func init() {
	//needed because init function from utils/kafka is run way before this file,
	//so the Cfg object is empty.
	Cfg.LoggingLevel = "DEBUG"
	//init the logger so it does not have to be initialized in each test
	setupLogger()
}

// TestHandleMessageProcessorNotSet verifies that the handleMessage panics if
// the processor is not set by the package instantiating the KafkaConsumer
func TestHandleMessageProcessorNotSet(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "message processor has not been set up. Aborting handleMessage", r)
		} else {
			t.Fail()
		}
	}()

	// construct dummy consumer without processor
	dummyConsumer := &KafkaConsumer{
		Config:                               &KafkaConsumerConfig{},
		ConsumerGroup:                        nil,
		numberOfSuccessfullyConsumedMessages: 0,
		numberOfErrorsConsumingMessages:      0,
		Ready:                                make(chan bool),
		Processor:                            nil,
	}

	message := sarama.ConsumerMessage{}
	message.Value = []byte(`{
		"OrgID": 5
	}`)

	dummyConsumer.handleMessage(&message)
}

func TestHandleMessageCheckCounters(t *testing.T) {
	// construct test processor
	processor := TestProcessor{}
	// construct dummy consumer, cannot use NewKafkaConsumer without setting
	// the correct env vars and having a running broker behing
	dummyConsumer := &KafkaConsumer{
		Config:                               &KafkaConsumerConfig{},
		ConsumerGroup:                        nil,
		numberOfSuccessfullyConsumedMessages: 0,
		numberOfErrorsConsumingMessages:      0,
		Ready:                                make(chan bool),
		Processor:                            &processor,
	}

	// check initial counters
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), processor.GetNumberOfMessagesProcessed())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message := sarama.ConsumerMessage{}
	message.Value = []byte(`{
		"OrgID": 5
	}`)

	// process message then check counters
	dummyConsumer.handleMessage(nil)
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), processor.GetNumberOfMessagesProcessed())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	dummyConsumer.handleMessage(&message)
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(1), processor.GetNumberOfMessagesProcessed())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())
}

func TestSetupLoggerFail(t *testing.T) {
	prevLogFatalf := logFatalf
	prevLogger := logger
	prevLoggingLevel := Cfg.LoggingLevel
	Cfg.LoggingLevel = "invalid"

	defer func() {
		Cfg.LoggingLevel = prevLoggingLevel
		logger = prevLogger
		logFatalf = prevLogFatalf
	}()

	logFatalf = func(format string, args ...interface{}) {
		assert.Equal(t, "Error setting up logger: invalid loglevel given", fmt.Sprintf(format, args...))
		logger = prevLogger // Graceful return
	}

	logger = nil
	setupLogger()
}

func initKafkaBroker(sasl, address, consumerGroup, incomingTopic string) {
	Cfg.KafkaBroker = createTestBroker(sasl)
	Cfg.KafkaBrokerAddress = address
	Cfg.KafkaBrokerConsumerGroup = consumerGroup
	Cfg.KafkaBrokerIncomingTopic = incomingTopic
}

func TestNewKafkaConsumer(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "test-broker-addr", "test-consumer-group", "test-broker-inc-topic")

	cfg := sarama.NewConfig()
	cfg.Metadata.Full = false
	_, err := NewKafkaConsumer(cfg, nil)
	assert.Nil(t, err)
}

func TestNewKafkaConsumerInvalidAddress(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "", "", "")

	cfg := sarama.NewConfig()
	cfg.Metadata.Full = false
	_, err := NewKafkaConsumer(cfg, nil)
	assert.Equal(t, "unable to get env var: KAFKA_BROKER_ADDRESS", err.Error())
}

func TestNewKafkaConsumerInvalidConsumerGroup(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "test-broker-addr", "", "")

	cfg := sarama.NewConfig()
	cfg.Metadata.Full = false
	_, err := NewKafkaConsumer(cfg, nil)
	assert.Equal(t, "unable to get env var: KAFKA_BROKER_CONSUMER_GROUP", err.Error())
}

func TestNewKafkaConsumerInvalidTopic(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "test-broker-addr", "test-consumer-group", "")

	cfg := sarama.NewConfig()
	cfg.Metadata.Full = false
	_, err := NewKafkaConsumer(cfg, nil)
	assert.Equal(t, "unable to get env var: KAFKA_BROKER_INCOMING_TOPIC", err.Error())
}

func TestNewKafkaConsumerOutOfBrokers(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "test-broker-addr", "test-consumer-group", "test-broker-inc-topic")

	_, err := NewKafkaConsumer(sarama.NewConfig(), nil)
	assert.Equal(t, "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)", err.Error())
}

func TestNewKafkaConsumerNilConfig(t *testing.T) {
	initKafkaBroker(sarama.SASLTypePlaintext, "test-broker-addr", "test-consumer-group", "test-broker-inc-topic")

	_, err := NewKafkaConsumer(nil, nil)
	// Expected to fail on NewClient call with config.Metadata.Full of sarama set to true.
	assert.Equal(t, "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)", err.Error())
}

type testConsumer struct {
	done            chan bool
	consumeResponse error
	closeResponse   error
	cancelCtx       bool
	logInterceptor  *interceptor
}

func (c *testConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	if c.cancelCtx {
		handler.(*KafkaConsumer).Cancel()
	} else {
		c.done <- true
	}

	return c.consumeResponse
}

func (c testConsumer) Errors() <-chan error {
	return make(chan error)
}

func (c testConsumer) Close() error {
	return c.closeResponse
}

func (c testConsumer) Pause(partitions map[string][]int32) {}

func (c testConsumer) Resume(partitions map[string][]int32) {}

func (c testConsumer) PauseAll() {}

func (c testConsumer) ResumeAll() {}

func (c testConsumer) Serve() {}

func getConsumer() KafkaConsumer {
	consumer := KafkaConsumer{
		Config:                               &KafkaConsumerConfig{},
		ConsumerGroup:                        &testConsumer{done: make(chan bool, 1), logInterceptor: &interceptor{}},
		numberOfSuccessfullyConsumedMessages: 0,
		numberOfErrorsConsumingMessages:      0,
		Ready:                                nil,
		Cancel:                               nil,
		Processor:                            nil,
	}

	return consumer
}

func (c *testConsumer) await(t *testing.T) {
	select {
	case <-c.done:
		return
	case <-c.logInterceptor.done:
		return
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for consumer")
	}
}

func TestKafkaConsumerServe(t *testing.T) {
	consumer := getConsumer()
	go consumer.Serve()
	consumer.ConsumerGroup.(*testConsumer).await(t)
}

type interceptor struct {
	expectedMessage string
	expectedError   string
	done            chan bool
	t               *testing.T
}

func (i interceptor) Write(p []byte) (n int, err error) {
	logMsg := string(p[:])
	assert.True(i.t, strings.Contains(logMsg, i.expectedMessage))
	assert.True(i.t, strings.Contains(logMsg, i.expectedError))
	i.done <- true
	return 0, nil
}

func setupTestLogger(t *testing.T) {
	var err error
	logger, err = CreateLogger(Cfg.LoggingLevel)
	assert.Nil(t, err)
}

// Intercepts and asserts logrus messages then exits Serve loop or fails the test.
func TestKafkaConsumerServeFail(t *testing.T) {
	setupTestLogger(t)
	errMsg := "failed successfully"

	it := interceptor{
		expectedMessage: "Unable to recreate Kafka session",
		expectedError:   errMsg,
		done:            make(chan bool, 1),
		t:               t,
	}
	logger.AddHook(&writer.Hook{
		Writer:    it,
		LogLevels: []logrus.Level{logrus.ErrorLevel},
	})

	consumer := getConsumer()
	tc := consumer.ConsumerGroup.(*testConsumer)
	tc.consumeResponse = errors.New(errMsg)
	tc.logInterceptor = &it

	go consumer.Serve()

	tc.await(t)
}

// Intercepts and asserts logrus messages then exits Serve loop or fails the test.
func TestKafkaConsumerServeCancel(t *testing.T) {
	setupTestLogger(t)
	errMsg := "context canceled"

	it := interceptor{
		expectedMessage: "stopping consumer",
		expectedError:   errMsg,
		done:            make(chan bool, 1),
		t:               t,
	}
	logger.AddHook(&writer.Hook{
		Writer:    it,
		LogLevels: []logrus.Level{logrus.InfoLevel},
	})

	consumer := getConsumer()
	tc := consumer.ConsumerGroup.(*testConsumer)
	tc.logInterceptor = &it
	tc.cancelCtx = true

	go consumer.Serve()

	tc.await(t)
}

func TestKafkaConsumerClose(t *testing.T) {
	consumer := getConsumer()
	assert.Nil(t, consumer.Close())
}

func TestKafkaConsumerCloseWithCancel(t *testing.T) {
	consumer := getConsumer()
	ctx, cancel := context.WithCancel(context.Background())
	consumer.Cancel = cancel
	assert.Nil(t, consumer.Close())
	assert.NotNil(t, ctx.Err())
}

func TestKafkaConsumerCloseWithGroup(t *testing.T) {
	setupTestLogger(t)
	errMsg := "failed successfully"

	it := interceptor{
		expectedMessage: errMsg,
		expectedError:   "Unable to close consumer group",
		done:            make(chan bool, 1),
		t:               t,
	}
	logger.AddHook(&writer.Hook{
		Writer:    it,
		LogLevels: []logrus.Level{logrus.ErrorLevel},
	})

	consumer := getConsumer()
	tc := consumer.ConsumerGroup.(*testConsumer)
	tc.closeResponse = errors.New(errMsg)
	tc.logInterceptor = &it

	consumer.Close()
}

func TestKafkaConsumerIncrementErr(t *testing.T) {
	consumer := getConsumer()
	before := consumer.numberOfErrorsConsumingMessages
	consumer.IncrementNumberOfErrorsConsumingMessages()
	after := consumer.numberOfErrorsConsumingMessages
	assert.Equal(t, before+1, after)
}

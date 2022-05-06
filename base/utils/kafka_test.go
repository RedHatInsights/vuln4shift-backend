package utils

import (
	"errors"
	"testing"

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

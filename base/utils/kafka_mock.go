package utils

import (
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

type SaramaAsyncWriterMock struct {
	invalidKey           string
	ExpectedErrorMessage string
	WriteQueue           chan *sarama.ProducerMessage
	SuccessMessages      chan *sarama.ProducerMessage
	ErrorMessages        chan *sarama.ProducerError
	ProcessedMessages    []*sarama.ProducerMessage
	Done                 chan bool
}

func (tw *SaramaAsyncWriterMock) Input() chan<- *sarama.ProducerMessage {
	return tw.WriteQueue
}

func (tw *SaramaAsyncWriterMock) Successes() <-chan *sarama.ProducerMessage {
	return tw.SuccessMessages
}

func (tw *SaramaAsyncWriterMock) Errors() <-chan *sarama.ProducerError {
	return tw.ErrorMessages
}

func (tw *SaramaAsyncWriterMock) Close() error {
	tw.Done <- true
	return nil
}

func (tw *SaramaAsyncWriterMock) StartProcessing(t *testing.T) {
	// Start listening for incoming Kafka messages
	for {
		select {
		case <-tw.Done:
			return
		case msg := <-tw.WriteQueue:
			rawKey, err := msg.Key.Encode()
			assert.Nil(t, err)
			key := string(rawKey)

			if key == SaramaMockInvalidKey {
				tw.ErrorMessages <- &sarama.ProducerError{
					Msg: msg,
					Err: errors.New(tw.ExpectedErrorMessage),
				}
			} else {
				tw.SuccessMessages <- msg
				tw.ProcessedMessages = append(tw.ProcessedMessages, msg)
			}

		}
	}
}

const (
	SaramaMockInvalidKey = "test-invalid-key"
)

func CreateSaramaAsyncWriterMock() *SaramaAsyncWriterMock {
	return &SaramaAsyncWriterMock{
		invalidKey:      SaramaMockInvalidKey,
		WriteQueue:      make(chan *sarama.ProducerMessage),
		SuccessMessages: make(chan *sarama.ProducerMessage),
		ErrorMessages:   make(chan *sarama.ProducerError),
		Done:            make(chan bool),
	}
}

func CreateKafkaProducerMock(topic string, writer *SaramaAsyncWriterMock) *KafkaProducer {
	return &KafkaProducer{
		Config: &KafkaProducerConfig{
			Address: Cfg.KafkaBrokerAddress,
			Topic:   topic,
		},
		numberOfSuccessfullyProducedMessages: 0,
		numberOfErrorsProducingMessages:      0,
		Writer:                               writer,
		Enqueued:                             0,
	}
}

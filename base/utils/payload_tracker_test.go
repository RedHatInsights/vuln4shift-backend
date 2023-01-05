package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func awaitPayloadTrackerValidResponse(t *testing.T, kafkaProducer *KafkaProducer, timeout time.Duration, previousSuccessCount uint64) {
	for ts := time.Now(); time.Since(ts) < timeout; {
		if kafkaProducer.GetNumberOfSuccessfullyProducedMessages() > previousSuccessCount {
			assert.Equal(t, 0, kafkaProducer.Enqueued)
			return
		}
	}
	t.Fail()
}

func TestSendPayloadTrackerMessage(t *testing.T) {
	setupTestLogger(t)

	topic := "test-pt-topic"

	testWriter := CreateSaramaAsyncWriterMock()
	go testWriter.StartProcessing(t)

	testProducer := CreateKafkaProducerMock(topic, testWriter)
	defer testProducer.Close()

	ptEvent := NewPayloadTrackerEvent("test-key")

	assert.Nil(t, ptEvent.SendKafkaMessage(testProducer))
	awaitPayloadTrackerValidResponse(t, testProducer, time.Millisecond*500, 0)

	assert.Equal(t, uint64(1), testProducer.GetNumberOfSuccessfullyProducedMessages())
	assert.Equal(t, uint64(0), testProducer.GetNumberOfErrorsProducingMessages())
}

func TestPayloadTrackerStatusUpdate(t *testing.T) {
	ptEvent := NewPayloadTrackerEvent("key")

	ptEvent.UpdateStatusReceived()
	assert.Equal(t, PayloadTrackerStatusReceived, ptEvent.Status)

	ptEvent.UpdateStatusError("err")
	assert.Equal(t, PayloadTrackerStatusError, ptEvent.Status)
	assert.Equal(t, "err", ptEvent.StatusMsg)

	ptEvent.UpdateStatusSuccess()
	assert.Equal(t, PayloadTrackerStatusSuccess, ptEvent.Status)

	assert.NotNil(t, ptEvent.Date)
}

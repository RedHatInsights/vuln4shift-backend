package digestwriter_test

import (
	"app/digestwriter"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

// Unit test definitions for functions and methods defined in source file
// consumer.go

// NewDummyConsumer function returns a new, not running, instance of
// KafkaConsumer.
func NewDummyConsumer() *digestwriter.KafkaConsumer {
	digestwriter.SetupLogger()
	return &digestwriter.KafkaConsumer{
		Config:        digestwriter.KafkaConsumerConfig{},
		ConsumerGroup: nil,
		Storage:       nil,
		Ready:         nil,
		Cancel:        nil,
	}
}

// TestParseEmptyMessage checks how empty message is handled by
// consumer.
func TestParseEmptyMessage(t *testing.T) {
	// empty message to be parsed
	const emptyMessage = ""

	// try to parse the message
	_, err := digestwriter.ParseMessage([]byte(emptyMessage))

	// check for error - it should be reported
	assert.EqualError(t, err, "unexpected end of JSON input")
}

// TestParseMessageMissingRequiredFields checks how empty message is handled by
// the consumer.
func TestParseMessageMissingRequiredFields(t *testing.T) {
	// empty message to be parsed
	jsonMessage := []byte(`{
		"no_digests_msg": "true"
	}`)

	// try to parse the message
	_, err := digestwriter.ParseMessage([]byte(jsonMessage))

	// check for error - it should be reported
	assert.EqualError(t, err, "missing required attribute 'images'")
}

// TestParseMessageMissingRequiredFields checks how message with no digests
// is handled by the consumer.
func TestParseMessageNoDigests(t *testing.T) {
	// empty message to be parsed
	jsonMessage := []byte(`{
		"images": {}
	}`)

	// try to parse the message
	_, err := digestwriter.ParseMessage([]byte(jsonMessage))

	// check for error - it should be reported
	assert.EqualError(t, err, "received message does not contain any digest")
}

// TestExtractDigestsFromMessage verify extraction of digests from correct message
func TestExtractDigestsFromMessage(t *testing.T) {
	// message to be parsed
	jsonMessage := []byte(`{
		"any_other_field": "whatever",
		"images": {
			"first_digest": {
			  "extra_content": [
				"more_content_1",
				"more_content_2",
				"more_content_3"
			  ],
			  "extra_content": "extra_content_value"
			},
			"second_digest": {
			  "second_digest_inner_data": "some_value"
			},
			"third_digest": {
			  "extra_content": [
				"more_content_1",
				"more_content_2",
				"more_content_3"
			  ],
			  "extra_content": "extra_content_value",
			  "extra_content_2": "extra_content_2_value",
			  "extra_content_3": "extra_content_3_value"
			}
		},
		"some_other_field": [ 1, 2, 3, 4]
	}`)

	// try to parse the message
	parsed, err := digestwriter.ParseMessage(jsonMessage)
	assert.Nil(t, err, "JSON string should be parsed correctly")

	digests := digestwriter.ExtractDigestsFromMessage(parsed.Digests)
	assert.Equal(t, 3, len(digests))
	assert.Contains(t, digests, "first_digest")
	assert.Contains(t, digests, "second_digest")
	assert.Contains(t, digests, "third_digest")
}

// TestProcessEmptyMessage check the behaviour of function processMessage with
// empty message on input.
func TestProcessEmptyMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// prepare an empty message
	message := sarama.ConsumerMessage{}

	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check for errors - it should be reported
	assert.EqualError(t, err, "unexpected end of JSON input")
}

// TestProcessWrongMessageMissingFields check the behaviour of function processMessage when
// received message does not contain the 'images' field.
func TestProcessWrongMessageMissingFields(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// prepare a message with missing 'images' field
	message := sarama.ConsumerMessage{}
	// fill in a message payload
	ConsumerMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"
	}`

	message.Value = []byte(ConsumerMessage)
	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check for errors - it should be reported
	assert.EqualError(t, err, "missing required attribute 'images'")
}

// TestProcessWrongMessageEmptyImages check the behaviour of function processMessage when
// received message does not contain any digest.
func TestProcessWrongMessageEmptyImages(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// prepare a message with missing 'images' field
	message := sarama.ConsumerMessage{}
	// fill in a message payload
	ConsumerMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"images": {}
	}`

	message.Value = []byte(ConsumerMessage)
	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check for errors - it should be reported
	assert.EqualError(t, err, "received message does not contain any digest")
}

// TestProcessMessageWithExpectedFields check the behaviour of function ProcessMessage when
// received message contains the 'images' field.
func TestProcessMessageWithExpectedFields(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()
	storage, mock := NewMockStorage(t)
	dummyConsumer.Storage = storage

	patchCurrentTime()

	// expected SQL statements during this test
	expectedStatement := `INSERT INTO "image" ("modified_date","digest") VALUES ($1,$2) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs(time.Now().UTC(), "first_digest").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// prepare a message with missing 'images' field
	message := sarama.ConsumerMessage{}
	// fill in a message payload
	ConsumerMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"images": {
			"first_digest": {
			  "extra_content": [
				"more_content_1",
				"more_content_2",
				"more_content_3"
			  ],
			  "extra_content": "extra_content_value"
			}
		}
	}`

	message.Value = []byte(ConsumerMessage)
	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check no errors reported
	assert.Nil(t, err, "input message should be processed correctly")

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

func TestHandleMessageCheckCounters(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()
	storage, mock := NewMockStorage(t)
	dummyConsumer.Storage = storage

	patchCurrentTime()

	// expected SQL statements during this test
	expectedStatement := `INSERT INTO "image" ("modified_date","digest") VALUES ($1,$2) RETURNING "id"`
	mock.ExpectBegin()
	//mock.ExpectExec(expectedStatement).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedStatement)).
		WithArgs(time.Now().UTC(), "first_digest").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// check initial counters
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message := sarama.ConsumerMessage{}

	WrongInputMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"
	}`

	NoDigestsMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"images": {}
	}`

	CorrectMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"images": {
			"first_digest": {
			  "extra_content": [
				"more_content_1",
				"more_content_2",
				"more_content_3"
			  ],
			  "extra_content": "extra_content_value"
			}
		}
	}`

	message.Value = []byte(WrongInputMessage)
	// try to handle the message
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	// check counters after processing
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message.Value = []byte(NoDigestsMessage)
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(2), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message.Value = []byte(CorrectMessage)
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(2), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

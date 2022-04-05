package digestwriter_test

// Unit test definitions for functions and methods defined in source file
// consumer.go

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

func NewDummyConsumerWithStorage(t *testing.T) (*digestwriter.KafkaConsumer, sqlmock.Sqlmock) {
	digestwriter.SetupLogger()
	storage, mock := NewMockStorage(t)
	return &digestwriter.KafkaConsumer{
		Config:        digestwriter.KafkaConsumerConfig{},
		ConsumerGroup: nil,
		Storage:       storage,
		Ready:         nil,
		Cancel:        nil,
	}, mock
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

// TestParseMessageMissingRequiredFields checks for mandatory fields
func TestParseMessageMissingRequiredFields(t *testing.T) {
	// message to be parsed
	jsonMessage := []byte(`{
		"no_digests_msg": "true"
	}`)
	// try to parse the message
	_, err := digestwriter.ParseMessage(jsonMessage)
	// check for error - it should be reported
	assert.EqualError(t, err, "missing required attribute 'AccountNumber'")

	jsonMessage = []byte(`{
		"AccountNumber": 2
	}`)
	_, err = digestwriter.ParseMessage(jsonMessage)
	assert.EqualError(t, err, "missing required attribute 'OrgID'")

	jsonMessage = []byte(`{
		"AccountNumber": 2,
		"OrgID": 1
	}`)
	_, err = digestwriter.ParseMessage(jsonMessage)
	assert.EqualError(t, err, "missing required attribute 'ClusterName'")

	// message to be parsed
	jsonMessage = []byte(`{
		"ClusterName": "a_name"
	}`)
	_, err = digestwriter.ParseMessage(jsonMessage)
	assert.EqualError(t, err, "missing required attribute 'AccountNumber'")

	jsonMessage = []byte(`{
		"AccountNumber": 2,
		"OrgID": 1,
		"ClusterName": "a_name"
	}`)
	_, err = digestwriter.ParseMessage(jsonMessage)
	assert.EqualError(t, err, "missing required attribute 'Images'")
}

// TestParseMessageNoDigests checks that a valid message with no digests
// is handled successfully by the consumer.
func TestParseMessageNoDigests(t *testing.T) {
	digestwriter.SetupLogger()
	// message to be parsed
	jsonMessage := []byte(`{
		"AccountNumber": 2,
		"OrgID": 1,
		"ClusterName": "a_name",
		"Images": {
			"images": {}
		}
	}`)

	// try to parse the message
	parsed, err := digestwriter.ParseMessage(jsonMessage)
	// check that no errors occur
	assert.Nil(t, err, "parseMessage should not fail if it contains no digests")
	assert.NotNil(t, parsed.Images.Digests)
	assert.Equal(t, 0, len(*parsed.Images.Digests))
}

// TestExtractDigestsFromMessage verify extraction of digests from correct message
func TestExtractDigestsFromMessage(t *testing.T) {
	digestwriter.SetupLogger()
	// message to be parsed
	jsonMessage := []byte(`{
		"any_other_field": "whatever",
		"AccountNumber": 2,
		"OrgID": 1,
		"ClusterName": "a_name",
		"Images": {
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
			}
		},
		"some_other_field": [ 1, 2, 3, 4]
	}`)

	// try to parse the message
	parsed, err := digestwriter.ParseMessage(jsonMessage)
	assert.Nil(t, err, "JSON should have been parsed correctly")

	digests := digestwriter.ExtractDigestsFromMessage(parsed.Images.Digests)
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

	// prepare a message with a required filed missing
	message := sarama.ConsumerMessage{}
	ConsumerMessageNoAccount := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"
	}`
	ConsumerMessageNoOrgID := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"AccountNumber": 3
	}`
	ConsumerMessageNoClusterName := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"AccountNumber": 3,
		"OrgID": 2
	}`

	ConsumerMessageNoImages := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"AccountNumber": 3,
		"ClusterName": "test",
		"OrgID": 1
	}`
	// try to process the messages and check for errors
	message.Value = []byte(ConsumerMessageNoAccount)
	err := dummyConsumer.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'AccountNumber'")

	message.Value = []byte(ConsumerMessageNoOrgID)
	err = dummyConsumer.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'OrgID'")

	message.Value = []byte(ConsumerMessageNoClusterName)
	err = dummyConsumer.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'ClusterName'")

	message.Value = []byte(ConsumerMessageNoImages)
	err = dummyConsumer.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'Images'")
}

// TestProcessWrongMessageEmptyImages check the behaviour of function processMessage when
// received message does not contain any digest.
func TestProcessWrongMessageEmptyImages(t *testing.T) {
	// construct dummy consumer
	dummyConsumer, mock := NewDummyConsumerWithStorage(t)

	// prepare a message with empty 'images' field
	message := sarama.ConsumerMessage{}
	// fill in a message payload
	ConsumerMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
		"AccountNumber": 3,
		"ClusterName": "test",
		"OrgID": 4,
		"Images": {}
	}`

	message.Value = []byte(ConsumerMessage)
	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check for errors - nothing should be reported
	assert.Nil(t, err, "received message does not need to contain any digest")
	assert.Equal(t, 1, int(dummyConsumer.GetNumberOfMessagesWithEmptyDigests()))

	//check that no processMessage is aborted without any call to Storage
	assert.Nil(t, mock.ExpectationsWereMet(), "no SQL queries should have been made")
}

// expect these SQL statements to be called when consumed message is valid and has at least 1 digest
func setHappyPathExpectations(mock sqlmock.Sqlmock) {
	// expected SQL statements during this test [SIMPLIFIED. This behavior is tested in storage_test.go]
	expectedSelectFromAccount := `SELECT * FROM "account"`
	expectedInsertIntoAccount := `INSERT INTO "account"`
	expectedInsertIntoCluster := `INSERT INTO "cluster"`
	expectedSelectFromImage := `SELECT * FROM "image"`
	expectedInsertIntoClusterImage := `INSERT INTO "cluster_image"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelectFromAccount)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsertIntoAccount)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsertIntoCluster)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelectFromImage)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(expectedInsertIntoClusterImage)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

// TestProcessMessageWithExpectedFields check the behaviour of function ProcessMessage when
// received message contains the 'images' field.
func TestProcessMessageWithRequiredFields(t *testing.T) {
	// construct dummy consumer
	dummyConsumer, mock := NewDummyConsumerWithStorage(t)
	setHappyPathExpectations(mock)

	message := sarama.ConsumerMessage{}
	ConsumerMessage := `{
		"OrgID": 4,
		"AccountNumber": 3,
		"ClusterName": "84f7eedc-0000-0000-9d4d-000000000000",
		"Images": {
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
		}
	}`
	message.Value = []byte(ConsumerMessage)

	err := dummyConsumer.ProcessMessage(&message)
	assert.Nil(t, err, "input message should be processed correctly")
	// check  all SQL-related expectations were met
	checkAllExpectations(t, mock)
}

func TestHandleMessageCheckCounters(t *testing.T) {
	// construct dummy consumer
	dummyConsumer, mock := NewDummyConsumerWithStorage(t)
	setHappyPathExpectations(mock)

	// check initial counters
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfMessagesWithEmptyDigests())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message := sarama.ConsumerMessage{}

	WrongInputMessage := `{
		"pods": 1,
		"clusters": 2,
		"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"
	}`

	NoDigestsMessage := `{
		"OrgID": 4,
		"AccountNumber": 3,
		"ClusterName": "84f7eedc-0000-0000-9d4d-000000000000",
		"Images": {
			"pods": 1,
			"clusters": 2,
			"timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `",
			"images": {}
		}
	}`

	CorrectMessage := `{
		"OrgID": 5,
		"AccountNumber": 3,
		"ClusterName": "84f7eedc-0000-0000-9d4d-000000000000",
		"Images": {
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
		}
	}`

	// process message then check counters
	message.Value = []byte(WrongInputMessage)
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfMessagesWithEmptyDigests())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message.Value = []byte(NoDigestsMessage)
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfMessagesWithEmptyDigests())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	message.Value = []byte(CorrectMessage)
	digestwriter.HandleKafkaMessage(dummyConsumer, &message)
	assert.Equal(t, uint64(2), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfMessagesWithEmptyDigests())
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// check if all expectations were met
	checkAllExpectations(t, mock)
}

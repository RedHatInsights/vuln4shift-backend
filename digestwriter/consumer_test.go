package digestwriter_test

// Unit test definitions for functions and methods defined in source file
// consumer.go

import (
	"app/base/utils"
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

func init() {
	// needed because init function from utils/kafka is run way before this file,
	// so the Cfg object is empty.
	utils.Cfg.LoggingLevel = "DEBUG"
	// init the logger so it does not have to be initialized in each test
	digestwriter.SetupLogger()
}

// NewDummyConsumerWithProcessor function returns a new, not running, instance of
// KafkaConsumer as well as the Processor it uses.
func NewDummyConsumerWithProcessor(t *testing.T) (*utils.KafkaConsumer, *digestwriter.DigestConsumer, sqlmock.Sqlmock) {
	storage, mock := NewMockStorage(t)
	consumer, processor := digestwriter.NewDummyConsumerWithProcessor(storage)
	return consumer, processor, mock
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
	processor := digestwriter.DigestConsumer{}
	// prepare an empty message
	message := sarama.ConsumerMessage{}
	// try to process the message
	err := processor.ProcessMessage(&message)
	// check for errors - it should be reported
	assert.EqualError(t, err, "unexpected end of JSON input")
}

// TestProcessWrongMessageMissingFields check the behaviour of the ProcessMessage
// function when received message does not contain the required fields.
func TestProcessWrongMessageMissingFields(t *testing.T) {
	processor := digestwriter.DigestConsumer{}
	// prepare a message with a required field missing
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
	err := processor.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'AccountNumber'")

	message.Value = []byte(ConsumerMessageNoOrgID)
	err = processor.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'OrgID'")

	message.Value = []byte(ConsumerMessageNoClusterName)
	err = processor.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'ClusterName'")

	message.Value = []byte(ConsumerMessageNoImages)
	err = processor.ProcessMessage(&message)
	assert.EqualError(t, err, "missing required attribute 'Images'")
}

// TestProcessWrongMessageEmptyImages check the behaviour of the ProcessMessage function
// when received message does not contain any digest.
func TestProcessWrongMessageEmptyImages(t *testing.T) {
	// construct dummy consumer just to make sure the processor is correctly set
	dummyConsumer, processor, mock := NewDummyConsumerWithProcessor(t)
	// prepare a message with empty 'Images' field
	message := sarama.ConsumerMessage{}
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

	// try to process the message using the consumer's Processor pointer
	err := dummyConsumer.Processor.ProcessMessage(&message)

	// check for errors - nothing should be reported
	assert.Nil(t, err, "received message does not need to contain any digest")
	// check that the counters are incremented accordingly
	assert.Equal(t, 1, int(processor.GetNumberOfMessagesWithEmptyDigests()))

	//check that no processMessage is aborted without any call to Storage
	assert.Nil(t, mock.ExpectationsWereMet(), "no SQL queries should have been made")
}

// expect these SQL statements to be called when consumed message is valid and has at least 1 digest
func setHappyPathExpectations(mock sqlmock.Sqlmock) {
	// expected SQL statements during this test [SIMPLIFIED. This behavior is tested in storage_test.go]
	expectedSelectFromAccount := `SELECT * FROM "account" WHERE`
	expectedInsertIntoAccount := `INSERT INTO "account"`
	expectedSelectFromCluster := `SELECT "cluster"."id","cluster"."uuid","cluster"."account_id" FROM "cluster"`
	expectedInsertIntoCluster := `INSERT INTO "cluster"`
	expectedSelectFromImage := `SELECT * FROM "image"`
	expectedInsertIntoClusterImage := `INSERT INTO "cluster_image"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelectFromAccount)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
	mock.ExpectQuery(regexp.QuoteMeta(expectedInsertIntoAccount)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(expectedSelectFromCluster)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
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
	dummyConsumer, processor, mock := NewDummyConsumerWithProcessor(t)
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

	err := dummyConsumer.Processor.ProcessMessage(&message)

	// check for errors - nothing should be reported
	assert.Nil(t, err, "input message should be processed correctly")
	// check that the counters are incremented accordingly
	assert.Equal(t, 0, int(processor.GetNumberOfMessagesWithEmptyDigests()))
	// check  all SQL-related expectations were met
	checkAllExpectations(t, mock)
}

package digestwriter

import (
	"app/base/utils"
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

const (
	// key for organization ID in structure log messages
	orgKey = "orgID"
	// key for account in structured log messages
	accountKey = "account"
	// key for cluster in structured log messages
	clusterKey = "cluster"
	// key for cluster ID retrieved from DB in structured log messages
	clusterIDKey = "clusterID"
	// key for consumed message version in structured log messages
	versionKey = "version"
	// key for request ID in structured log messages
	requestIDKey = "requestID"
	// key for error message used in structured log messages
	errorKey = "error"
)

// OrgID data type represents organization ID.
type OrgID uint32

// AccountNumber data type represents account number for a given report.
type AccountNumber uint32

// ClusterName represents the external cluster UUID contained in the consumed message
type ClusterName string

// RequestID represents the unique payload identifier of input Kafka records
type RequestID string

// JSONContent represents the workload info contained in the consumed message
type JSONContent map[string]*json.RawMessage

// Image data structure is representation of Images JSON object
type Image struct {
	Pods       int          `json:"-"`
	ImageCount int          `json:"imageCount"`
	Digests    *JSONContent `json:"images"`
	Namespaces *JSONContent `json:"-"`
}

// IncomingMessage data structure is representation of message consumed from
// the configured topic
type IncomingMessage struct {
	Organization  *OrgID         `json:"OrgID"`
	AccountNumber *AccountNumber `json:"AccountNumber"`
	ClusterName   *ClusterName   `json:"ClusterName"`
	Images        *Image         `json:"Images"`
	LastChecked   string         `json:"-"`
	Version       uint8          `json:"Version"`
	RequestID     RequestID      `json:"RequestID"`
}

// Struct that must fulfill the Processor interface defined in utils/kafka.go
// It is specific to each service, so it can have any requried fields not
// defined in the original Consumer interface
type DigestConsumer struct {
	storage                          Storage
	numberOfMessagesWithEmptyDigests uint64
}

// NewConsumer constructs a new instance of Consumer interface
// specialized in consuming from SHA extractor's result topic
func NewConsumer(storage Storage) (*utils.KafkaConsumer, error) {
	setupLogger()
	processor := DigestConsumer{
		storage,
		0,
	}
	consumer, err := utils.NewKafkaConsumer(nil, &processor)
	if err != nil {
		return nil, err
	}
	return consumer, err
}

// GetNumberOfMessagesWithEmptyDigests returns number of messages
// where Images field was empty
func (d *DigestConsumer) GetNumberOfMessagesWithEmptyDigests() uint64 {
	return d.numberOfMessagesWithEmptyDigests
}

// IncrementNumberOfMessagesWithEmptyDigests increments number of consumed message with no digests
func (d *DigestConsumer) IncrementNumberOfMessagesWithEmptyDigests() {
	d.numberOfMessagesWithEmptyDigests++
}

// ProcessMessage processes an incoming message
func (d *DigestConsumer) ProcessMessage(msg *sarama.ConsumerMessage) error {
	// Step #1: parse the incoming message
	message, err := parseMessage(msg.Value)
	if err != nil {
		/* ParseIncomingMessageError.Inc() */
		return err
	}

	/* ParsedIncomingMessage.Inc() */

	if message.Images.Digests == nil || len(*message.Images.Digests) == 0 {
		logger.Infoln("no digests were retrieved from incoming message")
		d.IncrementNumberOfMessagesWithEmptyDigests()
		return nil
	}

	// Step #2: get digests into a slice of strings
	digests := extractDigestsFromMessage(message.Images.Digests)

	logger.Debugf("extracted digests: %d\n", len(digests))

	if message.Images.ImageCount != len(digests) {
		logger.Warnf("Expected number of digests: %d; Extracted digests: %d\n",
			message.Images.ImageCount, len(digests))
	}

	// Step #3: update tables with received info
	err = d.storage.WriteClusterInfo(
		message.ClusterName, message.AccountNumber, message.Organization, digests)
	if err != nil {
		logger.WithFields(logrus.Fields{
			accountKey: message.AccountNumber,
			clusterKey: message.ClusterName,
			errorKey:   err.Error(),
		}).Errorln("error writing to cluster table")
		/* StoredMessagesError.Inc() */
		return err
	}

	/* StoredMessagesOk.Inc() */
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
	if deserialized.AccountNumber == nil {
		return deserialized, errors.New("missing required attribute 'AccountNumber'")
	}
	if deserialized.Organization == nil {
		return deserialized, errors.New("missing required attribute 'OrgID'")
	}
	if deserialized.ClusterName == nil {
		return deserialized, errors.New("missing required attribute 'ClusterName'")
	}
	if deserialized.Images == nil {
		return deserialized, errors.New("missing required attribute 'Images'")
	}

	logger.WithFields(logrus.Fields{
		requestIDKey: deserialized.RequestID,
		versionKey:   deserialized.Version,
		orgKey:       deserialized.Organization,
		accountKey:   deserialized.AccountNumber,
		clusterKey:   deserialized.ClusterName,
	}).Debugln("parsed incoming message correctly")

	return deserialized, nil
}

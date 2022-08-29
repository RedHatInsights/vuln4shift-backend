package digestwriter

import (
	"app/base/utils"
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

const (
	// Keys used in structured log messages
	// key for organization ID
	orgKey = "org_id"
	// key for account
	accountKey = "account"
	// key for cluster
	clusterKey = "cluster"
	// key for cluster ID
	clusterIDKey = "cluster_id"
	// key for consumed message version
	versionKey = "version"
	// key for request ID
	requestIDKey = "request_id"
	// key for error message
	errorKey = "error"
	// key for created row's ID in database table
	rowIDKey = "row_id"
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

type Container struct {
	ImageID string `json:"imageID"`
}

type Shape struct {
	InitContainers []Container `json:"initContainers"`
	Containers     []Container `json:"containers"`
}

// Workload data structure is representation of Images JSON object
type Workload struct {
	Pods       int                   `json:"-"`
	ImageCount int                   `json:"imageCount"`
	Images     *JSONContent          `json:"images"`
	Namespaces *map[string]Namespace `json:"namespaces"`
}

type Namespace struct {
	Shapes []Shape `json:"shapes"`
}

// IncomingMessage data structure is representation of message consumed from
// the configured topic
type IncomingMessage struct {
	Organization  OrgID         `json:"OrgID"`
	AccountNumber AccountNumber `json:"AccountNumber"`
	ClusterName   ClusterName   `json:"ClusterName"`
	Workload      *Workload     `json:"Images"`
	LastChecked   string        `json:"-"`
	Version       uint8         `json:"Version"`
	RequestID     RequestID     `json:"RequestID"`
}

// DigestConsumer Struct that must fulfill the Processor interface defined in utils/kafka.go
// It is specific to each service, so it can have any required fields not
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
		parseIncomingMessageError.Inc()
		return err
	}

	parsedIncomingMessage.Inc()

	if message.Workload.Images == nil || message.Workload.Namespaces == nil {
		logger.Infoln("no digests were retrieved from incoming message")
		d.IncrementNumberOfMessagesWithEmptyDigests()
		return nil
	}

	// Step #2: get digests into a slice of strings
	digests := extractDigestsFromMessage(*message.Workload)

	logger.Debugf("number of extracted digests: %d", len(digests))

	if message.Workload.ImageCount != len(digests) {
		logger.Warnf("expected number of digests: %d, extracted digests: %d",
			message.Workload.ImageCount, len(digests))
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
		storedMessagesError.Inc()
		return err
	}

	storedMessagesOk.Inc()
	return nil
}

func extractDigestsFromMessage(workload Workload) (digests []string) {
	digestSet := map[string]struct{}{}
	for imageID := range *workload.Images {
		digestSet[imageID] = struct{}{}
	}
	for _, namespace := range *workload.Namespaces {
		for _, shape := range namespace.Shapes {
			for _, initContainer := range shape.InitContainers {
				digestSet[initContainer.ImageID] = struct{}{}
			}
			for _, container := range shape.Containers {
				digestSet[container.ImageID] = struct{}{}
			}
		}
	}
	digests = []string{}
	for imageID := range digestSet {
		digests = append(digests, imageID)
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

	logger.Debugf("deserialized message: %v", deserialized)

	if deserialized.Workload == nil {
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

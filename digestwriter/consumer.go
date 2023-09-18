package digestwriter

import (
	"app/base/utils"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	usePayloadTracker bool
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

	// first byte of gzip compression
	gzip1 = 31
	// second byte of gzip compression
	gzip2 = 139
)

const (
	errNoDigests   = "no digests were retrieved from incoming message"
	errClusterData = "error updating cluster data"
)

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
	Organization  json.Number `json:"OrgID"`
	AccountNumber json.Number `json:"AccountNumber"`
	ClusterName   ClusterName `json:"ClusterName"`
	Workload      *Workload   `json:"Images"`
	LastChecked   string      `json:"-"`
	Version       uint8       `json:"Version"`
	RequestID     RequestID   `json:"RequestID"`
}

// DigestConsumer Struct that must fulfill the Processor interface defined in utils/kafka.go
// It is specific to each service, so it can have any required fields not
// defined in the original Consumer interface
type DigestConsumer struct {
	storage                          Storage
	numberOfMessagesWithEmptyDigests uint64
	PayloadTracker                   *utils.KafkaProducer
}

// startPayloadTracker starts Payload Tracker Kafka producer.
func startPayloadTracker() (*utils.KafkaProducer, error) {
	ptWriter, err := utils.NewKafkaProducer(nil, utils.Cfg.KafkaBrokerAddress, utils.Cfg.KafkaPayloadTrackerTopic)
	if err != nil {
		return nil, err
	}

	return ptWriter, nil
}

// NewConsumer constructs a new instance of Consumer interface
// specialized in consuming from SHA extractor's result topic
func NewConsumer(storage Storage) (*utils.KafkaConsumer, error) {
	SetupLogger()
	usePayloadTracker = utils.Cfg.PayloadTrackerEnabled

	processor := DigestConsumer{
		storage,
		0,
		nil,
	}

	if usePayloadTracker {
		payloadTracker, err := startPayloadTracker()
		if err != nil {
			return nil, err
		}
		processor.PayloadTracker = payloadTracker
	}

	consumer, err := utils.NewKafkaConsumer(nil, &processor)
	if err != nil {
		if usePayloadTracker {
			processor.PayloadTracker.Close()
		}
		return nil, err
	}

	if usePayloadTracker {
		// Release Payload Tracker producer resources during DigestConsumer Close
		consumer.Shutdown = processor.PayloadTracker.Close
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
	logger.Debugf("processing incoming message with a key=%s", msg.Key)

	// Step #1: parse the incoming message
	message, err := parseMessage(msg.Value)
	if err != nil {
		parseIncomingMessageError.Inc()
		return err
	}
	parsedIncomingMessage.Inc()

	// Set up payload tracker event
	ptEvent := utils.NewPayloadTrackerEvent(string(message.RequestID))
	ptEvent.SetOrgIDFromUint(message.Organization)

	// Send Payload Tracker message with status received
	ptEvent.UpdateStatusReceived()
	go d.sendPayloadTrackerMessage(ptEvent)

	if message.Workload.Images == nil || message.Workload.Namespaces == nil {
		logger.Debugln(errNoDigests)
		d.IncrementNumberOfMessagesWithEmptyDigests()
		ptEvent.UpdateStatusError(errNoDigests)
		go d.sendPayloadTrackerMessage(ptEvent)
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
		message.ClusterName, message.Organization, *message.Workload, digests)
	if err != nil {
		logger.WithFields(logrus.Fields{
			orgKey:     message.Organization,
			clusterKey: message.ClusterName,
			errorKey:   err.Error(),
		}).Errorln(errClusterData)
		storedMessagesError.Inc()
		ptEvent.UpdateStatusError(errClusterData)
		go d.sendPayloadTrackerMessage(ptEvent)
		return err
	}

	ptEvent.UpdateStatusSuccess()
	go d.sendPayloadTrackerMessage(ptEvent)
	storedMessagesOk.Inc()
	return nil
}

// sendPayloadTrackerMessage sends Kafka message to Payload Tracker and logs errors
func (d *DigestConsumer) sendPayloadTrackerMessage(event utils.PayloadTrackerEvent) {
	if !usePayloadTracker {
		return
	}

	logger.Debugf("sending Payload Tracker message with status %s", event.Status)
	if err := event.SendKafkaMessage(d.PayloadTracker); err != nil {
		logger.Errorf("failed to send Payload Tracker message: %s", err.Error())
		payloadTrackerError.Inc()
		return
	}

	payloadTrackerMessageSent.Inc()
	logger.Debugf("successfully sent Payload Tracker message req=%v, ts=%s, status=%s", *event.RequestID, *event.Date, event.Status)
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

// isMessageInGzipFormat will check if the format of the message is gzip if it is it will return true if not it will return false
func isMessageInGzipFormat(messageValue []byte) bool {
	if messageValue == nil {
		return false
	}
	if len(messageValue) < 2 {
		return false
	}
	// Checking for first 2 bytes in gzip instance witch are 31 and 139
	if messageValue[0] == gzip1 && messageValue[1] == gzip2 {
		return true
	}
	return false

}

// decompressMessage will try to decompress the message if the message is compressed
func decompressMessage(messageValue []byte) ([]byte, error) {
	if isMessageInGzipFormat(messageValue) {
		reader := bytes.NewReader(messageValue)
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer func(r *gzip.Reader) {
			if err := r.Close(); err != nil {
				logger.Errorf("failed to close gzip reader: %s", err.Error())
			}
		}(gzipReader)
		decompresed, err := io.ReadAll(gzipReader)
		if err != nil {
			return nil, err
		}
		return decompresed, err
	}
	return messageValue, nil
}

// parseMessage tries to parse incoming message and verify all required attributes
func parseMessage(messageValue []byte) (IncomingMessage, error) {
	var deserialized IncomingMessage
	messageValue, err := decompressMessage(messageValue)
	if err != nil {
		return deserialized, errors.Wrap(err, "failed to decompress incoming message")
	}
	err = json.Unmarshal(messageValue, &deserialized)

	if err != nil {
		return deserialized, err
	}

	logger.Debugf("deserialized message: %v", deserialized)

	if deserialized.Workload == nil {
		return deserialized, errors.New("missing required attribute 'Images'")
	}
	if deserialized.Organization.String() == "" {
		return deserialized, errors.New("OrgID cannot be null")
	}

	logger.WithFields(logrus.Fields{
		requestIDKey: deserialized.RequestID,
		versionKey:   deserialized.Version,
		orgKey:       deserialized.Organization,
		clusterKey:   deserialized.ClusterName,
	}).Debugln("parsed incoming message correctly")

	return deserialized, nil
}

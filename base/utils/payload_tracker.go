package utils

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

const (
	service = "ocp-vulnerability"

	// Payload Tracker required statuses
	PayloadTrackerStatusReceived = "received"
	PayloadTrackerStatusSuccess  = "success"
	PayloadTrackerStatusError    = "error"
)

type PayloadTrackerEvent struct {
	Service     string  `json:"service"`
	OrgID       *string `json:"org_id,omitempty"`
	RequestID   *string `json:"request_id"`
	InventoryID string  `json:"inventory_id"`
	Status      string  `json:"status"`
	StatusMsg   string  `json:"status_msg,omitempty"`
	Date        *string `json:"date"` // RFC3339
}

func (e *PayloadTrackerEvent) SetOrgIDFromUint(id uint32) {
	orgID := strconv.FormatUint(uint64(id), 10)
	e.OrgID = &orgID
}

func (e *PayloadTrackerEvent) SetRequestID(id string) {
	e.RequestID = &id
}

func (e *PayloadTrackerEvent) UpdateStatusReceived() {
	e.Status = PayloadTrackerStatusReceived
	e.updateTimestamp()
}

func (e *PayloadTrackerEvent) UpdateStatusSuccess() {
	e.Status = PayloadTrackerStatusSuccess
	e.updateTimestamp()
}

func (e *PayloadTrackerEvent) UpdateStatusError(msg string) {
	e.Status = PayloadTrackerStatusError
	e.StatusMsg = msg
	e.updateTimestamp()
}

func (e *PayloadTrackerEvent) updateTimestamp() {
	ts := time.Now().Format(time.RFC3339)
	e.Date = &ts
}

// SendKafkaMessage delegates sending Kafka message to Producer in non-blocking manner
func (e *PayloadTrackerEvent) SendKafkaMessage(producer Producer) error {
	bs, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(err, "failed to marshal Payload Tracker event")
	}

	go producer.SendMessage(sarama.StringEncoder(*e.RequestID), sarama.ByteEncoder(bs))

	return nil
}

func NewPayloadTrackerEvent(reqID string) PayloadTrackerEvent {
	return PayloadTrackerEvent{
		Service:     service,
		OrgID:       nil,
		RequestID:   &reqID,
		InventoryID: "",
		Status:      "",
		StatusMsg:   "",
		Date:        nil,
	}
}

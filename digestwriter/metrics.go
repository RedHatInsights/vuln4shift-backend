package digestwriter

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	parsedIncomingMessage = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of parsed incoming messages",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "parsed_incoming_messages",
	})

	parseIncomingMessageError = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of parse incoming message errors",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "parse_incoming_messages_error",
	})

	storedMessagesOk = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of successfully stored messages",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "stored_messages_ok",
	})

	storedMessagesError = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of unsuccessfully stored messages",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "stored_messages_error",
	})

	payloadTrackerMessageSent = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of messages sent to Payload Tracker",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "payload_tracker_message_sent",
	})

	payloadTrackerError = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of Payload Tracker errors",
		Namespace: "vuln4shift",
		Subsystem: "digestwriter",
		Name:      "payload_tracker_error",
	})
)

func RunMetrics() {
	prometheus.MustRegister(
		parsedIncomingMessage,
		parseIncomingMessageError,
		storedMessagesOk,
		storedMessagesError,
		payloadTrackerMessageSent,
		payloadTrackerError,
		utils.ConsumedMessages,
		utils.ConsumingErrors,
	)
	utils.StartPrometheus("vuln4shift_digestwriter")
}

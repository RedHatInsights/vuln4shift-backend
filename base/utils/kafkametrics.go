package utils

import "github.com/prometheus/client_golang/prometheus"

var (
	ConsumingErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of unsuccessfully consumed messages",
		Namespace: "vuln4shift",
		Subsystem: "kafka",
		Name:      "consuming_errors",
	})

	ConsumedMessages = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "Number of successfully consumed messages",
		Namespace: "vuln4shift",
		Subsystem: "kafka",
		Name:      "consumed_messages",
	})
)

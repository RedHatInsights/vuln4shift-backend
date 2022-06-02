package vmsync

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	dbConnection   = "db-connection"
	dbFetch        = "db-fetch"
	dbInsertUpdate = "db-insert-update"
	/* dbDelete     = "db-delete" */
)

var (
	syncError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many errors and of which type",
		Namespace: "vuln4shift",
		Subsystem: "vmsync",
		Name:      "sync_error",
	}, []string{"type"})

	vmaasRequestError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "Vmaas api request error by status code",
		Namespace: "vuln4shift",
		Subsystem: "vmsync",
		Name:      "request_error",
	}, []string{"url", "method", "code"})

	cvesInsertedUpdated = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many cves were inserted/updated during sync with VMAAS",
		Namespace: "vuln4shift",
		Subsystem: "vmsync",
		Name:      "cves_synced",
	})

	/* cvesDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many cves were deleted during sync with VMAAS",
		Namespace: "vuln4shift",
		Subsystem: "vmsync",
		Name:      "cves_synced",
	}) */
)

func GetMetricsPusher() *push.Pusher {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		syncError,
		vmaasRequestError,
		cvesInsertedUpdated,
	)
	pusher := push.New("http://pushgateway:9091", "vmsync").Gatherer(registry)

	return pusher
}

package vmsync

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	dbConnection   = "db-connection"
	dbFetch        = "db-fetch"
	dbInsertUpdate = "db-insert-update"
	dbDelete       = "db-delete"

	job = "vmsync"
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

	cvesDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many CVEs were deleted during sync with VMaaS",
		Namespace: "vuln4shift",
		Subsystem: "vmsync",
		Name:      "cves_deleted",
	})
)

func getMetricsPusher() *push.Pusher {
	return utils.GetMetricsPusher(
		job,
		syncError,
		vmaasRequestError,
		cvesInsertedUpdated,
		cvesDeleted)
}

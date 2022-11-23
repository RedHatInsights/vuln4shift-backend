package expsync

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	dbConnection   = "db-connection"
	dbInsertUpdate = "db-insert-update"
	dbDelete       = "db-delete"

	job       = "expsync"
	namespace = "vuln4shift"
)

var (
	syncError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many errors and of which type",
		Namespace: namespace,
		Subsystem: job,
		Name:      "sync_error",
	}, []string{"type"})

	exploitsRequestError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "exploit file request error by status code",
		Namespace: namespace,
		Subsystem: job,
		Name:      "request_error",
	}, []string{"url", "method", "code"})

	cveExploitsInsertedUpdated = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many CVEs have had inserted/updated exploits during sync with ProdSec API",
		Namespace: namespace,
		Subsystem: job,
		Name:      "cves_exploits_synced",
	})

	cveExploitsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many CVEs have had deleted exploits during sync with ProdSec API",
		Namespace: namespace,
		Subsystem: job,
		Name:      "cves_exploits_deleted",
	})
)

func getMetricsPusher() *push.Pusher {
	return utils.GetMetricsPusher(
		job,
		syncError,
		exploitsRequestError,
		cveExploitsInsertedUpdated,
		cveExploitsDeleted)
}

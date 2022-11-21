package expsync

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	dbConnection = "db-connection"
)

var (
	syncError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many errors and of which type",
		Namespace: "vuln4shift",
		Subsystem: "expsync",
		Name:      "sync_error",
	}, []string{"type"})

	exploitsRequestError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "exploit file request error by status code",
		Namespace: "vuln4shift",
		Subsystem: "expsync",
		Name:      "request_error",
	}, []string{"url", "method", "code"})
)

func getMetricsPusher() *push.Pusher {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		syncError,
	)
	pusher := push.New(utils.Cfg.PrometheusPushGateway, "expsync").Gatherer(registry)

	return pusher
}

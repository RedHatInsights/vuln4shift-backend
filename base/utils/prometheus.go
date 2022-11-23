package utils

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var logFatalf = log.Fatalf

func exposeOnPort(app *gin.Engine, port int) {
	err := app.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logFatalf(err.Error())
	}
}

func StartPrometheus(subsystem string) {
	app := gin.New()
	metricsPort := Cfg.MetricsPort

	p := ginprometheus.NewPrometheus(subsystem)
	p.MetricsPath = Cfg.MetricsPath
	p.Use(app)

	if metricsPort == -1 {
		go exposeOnPort(app, Cfg.PublicPort)
	} else {
		go exposeOnPort(app, metricsPort)
	}
}

func GetMetricsPusher(job string, collectors ...prometheus.Collector) *push.Pusher {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors...)

	return push.New(Cfg.PrometheusPushGateway, job).Gatherer(registry)
}

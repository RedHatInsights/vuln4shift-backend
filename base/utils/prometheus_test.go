package utils

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func checkPrometheus(t *testing.T, port int) {
	resp, err := http.Get(fmt.Sprintf("http://:%d%s", port, Cfg.MetricsPath))
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPrometheus(t *testing.T) {
	StartPrometheus("test_subsystem")
	time.Sleep(time.Millisecond * 50)

	port := Cfg.MetricsPort
	if port == -1 {
		port = Cfg.PublicPort
	}
	checkPrometheus(t, port)
}

func TestPrometheusPublic(t *testing.T) {
	prevPublicPort := Cfg.PublicPort
	Cfg.PublicPort = 7357
	prevMetricsPort := Cfg.MetricsPort
	Cfg.MetricsPort = -1

	defer func() {
		Cfg.MetricsPort = prevMetricsPort
		Cfg.PublicPort = prevPublicPort
	}()

	StartPrometheus("test_subsystem")
	time.Sleep(time.Millisecond * 50)

	checkPrometheus(t, Cfg.PublicPort)
}

func TestExposeOnPortFatal(t *testing.T) {
	prevLogFatalf := logFatalf
	defer func() { logFatalf = prevLogFatalf }()

	portOverflow := math.MaxUint16 + 1
	logFatalf = func(format string, args ...interface{}) {
		assert.Equal(t, fmt.Sprintf("listen tcp: address %d: invalid port", portOverflow), fmt.Sprintf(format, args...))
		return
	}

	exposeOnPort(gin.New(), portOverflow)
}

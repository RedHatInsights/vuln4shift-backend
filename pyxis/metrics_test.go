package pyxis

import (
	"app/base/utils"
	"app/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricsPusher(t *testing.T) {
	srv := test.GetMetricsServer(t, "PUT", "pyxis")
	defer srv.Close()

	oldPrometheusGateway := utils.Cfg.PrometheusPushGateway
	defer func() { utils.Cfg.PrometheusPushGateway = oldPrometheusGateway }()
	utils.Cfg.PrometheusPushGateway = srv.URL

	pusher := GetMetricsPusher()
	assert.Nil(t, pusher.Push())
}

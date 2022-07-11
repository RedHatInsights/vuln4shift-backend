package base

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus"
)

func RunMetrics() {
	prometheus.MustRegister()
	utils.StartPrometheus("vuln4shift_manager")
}

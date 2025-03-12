package measurements

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var measurements = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "count",
}, []string{"operation", "status"})

func UpdateCount(operation string, err error) {
	status := "success"
	if err != nil {
		status = "failure"
	}
	measurements.With(prometheus.Labels{"operation": operation, "status": status}).Inc()
}

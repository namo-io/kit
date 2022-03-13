package metric

import (
	"fmt"
	"net/http"

	"github.com/namo-io/kit/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	defaultLabels prometheus.Labels
)

func SetPrometheusMetricProvider(serviceName string, serviceId string, serviceVersion string, metricPort int) error {
	if len(serviceName) == 0 {
		return fmt.Errorf("service name is empty")
	}

	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%v", metricPort), nil)
	}()

	defaultLabels = prometheus.Labels{
		"service.name":    serviceName,
		"service.id":      serviceId,
		"service.version": serviceVersion,
		"host":            util.GetHostname(),
	}

	return nil
}

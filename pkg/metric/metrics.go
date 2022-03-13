package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GrpcRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "grpc_request_total",
		Help:        "grpc requests total",
		ConstLabels: defaultLabels,
	}, []string{"method", "status_code", "exception_type"})

	GrpcRequestDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "grpc_request_duration_seconds",
		Help:        "grpc requests latency per seconds",
		ConstLabels: defaultLabels,
	}, []string{"method", "status_code", "exception_type"})
)

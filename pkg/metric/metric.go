package metric

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	mp           = metric.NewNoopMeterProvider()
	meter        = mp.Meter("")
	defaultAttrs = []attribute.KeyValue{}
)

func Meter() metric.Meter {
	return meter
}

func NewAttrbiutesWithDefaults(attrs ...attribute.KeyValue) []attribute.KeyValue {
	return append(defaultAttrs, attrs...)
}

func SetPrometheusMetricProvider(serviceName string, serviceId string, serviceVersion string, metricPort int) error {
	if len(serviceName) == 0 {
		return fmt.Errorf("service name is empty")
	}

	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	exporter, err := prometheus.New(config, c)
	if err != nil {
		return fmt.Errorf("failed to initialize prometheus exporter %v", err)
	}

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(":8801", nil)
	}()

	mp = exporter.MeterProvider()
	meter = mp.Meter("")

	defaultAttrs = []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceInstanceIDKey.String(serviceId),
		semconv.ServiceVersionKey.String(serviceVersion),
	}

	return nil
}

package trace

import (
	"context"
	"fmt"

	"github.com/namo-io/kit/pkg/mctx"
	"github.com/namo-io/kit/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	otrace "go.opentelemetry.io/otel/trace"
)

var (
	defaultTp = tracesdk.NewTracerProvider()
	gtp       = defaultTp
)

func Start(ctx context.Context, spanName string) (context.Context, otrace.Span) {
	attrs := []attribute.KeyValue{}

	requestId := mctx.GetRequestId(ctx)
	if len(requestId) != 0 {
		attrs = append(attrs, attribute.String("request.id", requestId))
	}

	return gtp.Tracer("").Start(ctx, spanName, otrace.WithAttributes(attrs...))
}

func Shutdown(ctx context.Context) error {
	if defaultTp == gtp {
		return nil
	}

	return gtp.Shutdown(ctx)
}

func SetJeagerTraceProvider(serviceName string, serviceId string, serviceVersion string, jeagerEndpoint string) error {
	if len(serviceName) == 0 {
		return fmt.Errorf("service name is empty")
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jeagerEndpoint)))
	if err != nil {
		return err
	}

	attrs := []attribute.KeyValue{
		attribute.String("host", util.GetHostname()),
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceInstanceIDKey.String(serviceId),
		semconv.ServiceVersionKey.String(serviceVersion),
	}

	gtp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(semconv.SchemaURL, attrs...),
		),
	)

	return nil
}

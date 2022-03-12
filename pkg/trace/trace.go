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
	gtp = tracesdk.NewTracerProvider()
)

func Start(ctx context.Context, spanName string) (context.Context, otrace.Span) {
	attrs := []attribute.KeyValue{}

	requestId := mctx.GetRequestId(ctx)
	if len(requestId) == 0 {
		attrs = append(attrs, attribute.String("request.id", requestId))
	}

	return gtp.Tracer("").Start(ctx, spanName, otrace.WithAttributes(attrs...))
}

func Shutdown(ctx context.Context) error {
	return gtp.Shutdown(ctx)
}

func SetJeagerTraceProviderFromContext(ctx context.Context, url string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}

	attrs := []attribute.KeyValue{
		attribute.String("host", util.GetHostname()),
	}

	appName := mctx.GetAppName(ctx)
	if len(appName) == 0 {
		return fmt.Errorf("app name is nil")
	}
	attrs = append(attrs, semconv.ServiceNameKey.String(appName))

	appId := mctx.GetAppId(ctx)
	if len(appId) != 0 {
		attrs = append(attrs, semconv.ServiceInstanceIDKey.String(appId))
	}

	appVersion := mctx.GetAppVersion(ctx)
	if len(appVersion) != 0 {
		attrs = append(attrs, semconv.ServiceVersionKey.String(appVersion))
	}

	gtp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(semconv.SchemaURL, attrs...),
		),
	)

	return nil
}

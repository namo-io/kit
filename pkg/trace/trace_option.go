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
)

type TraceOption func(*Tracer) error

func WithJeagerTraceProviderFromContext(ctx context.Context, url string) TraceOption {
	return func(t *Tracer) error {
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

		t.tp = tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exp),
			tracesdk.WithResource(
				resource.NewWithAttributes(semconv.SchemaURL, attrs...),
			),
		)
		return nil
	}
}

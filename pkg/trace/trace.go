package trace

import (
	"context"

	"github.com/namo-io/kit/pkg/mctx"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	otrace "go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tp *tracesdk.TracerProvider
}

func NewTracer(opts ...TraceOption) (*Tracer, error) {
	t := &Tracer{
		tp: tracesdk.NewTracerProvider(),
	}

	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *Tracer) Start(ctx context.Context, spanName string) (context.Context, otrace.Span) {
	attrs := []attribute.KeyValue{}

	requestId := mctx.GetRequestId(ctx)
	if len(requestId) == 0 {
		attrs = append(attrs, attribute.String("request.id", requestId))
	}

	return t.tp.Tracer("").Start(ctx, spanName, otrace.WithAttributes(attrs...))
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	return t.tp.Shutdown(ctx)
}

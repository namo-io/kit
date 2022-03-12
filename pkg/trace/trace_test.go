package trace

import (
	"context"
	"testing"
	"time"

	"github.com/namo-io/kit/pkg/mctx"
)

func Test(t *testing.T) {
	ctx := mctx.WithAppName(context.Background(), "tracesdk-test")

	tc, err := NewTracer(WithJeagerTraceProviderFromContext(ctx, "http://localhost:14268/api/traces"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, span := tc.Start(ctx, "Qqwe")
	span.End()

	if err := tc.tp.ForceFlush(ctx); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	if err := tc.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
}

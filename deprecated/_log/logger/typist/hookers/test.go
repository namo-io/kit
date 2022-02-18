package hooks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/namo-io/kit/pkg/log/logger/typist"
)

type testHooker struct {
	errorMode bool
}

func (t *testHooker) Name() string {
	return "TestHooker"
}

func (t *testHooker) Fire(ctx context.Context, level typist.Level, rs *typist.Record) error {
	if t.errorMode {
		ctx, _ = context.WithTimeout(ctx, time.Second*3)
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			return err
		}

		return errors.New("errrorororororororor!!!\r")
	}

	fmt.Println("TESTSETSETSE Hooker")
	return nil
}

func NewTestHooker(errorMode bool) typist.Hooker {
	return &testHooker{
		errorMode: errorMode,
	}
}

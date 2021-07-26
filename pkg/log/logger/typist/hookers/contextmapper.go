package hooks

import (
	"context"

	"github.com/namo-io/kit/pkg/keys"
	"github.com/namo-io/kit/pkg/log/logger/typist"
)

type contextMapper struct {
}

func NewContextMapper() *contextMapper {
	return &contextMapper{}
}

func (t *contextMapper) Name() string {
	return "ContextMapper"
}

func (c *contextMapper) Fire(ctx context.Context, level typist.Level, rs *typist.Record) error {
	if contextRequiestID := ctx.Value(keys.RequestID); contextRequiestID != nil {
		rs.Meta[keys.RequestID.String()] = contextRequiestID
	}

	return nil
}

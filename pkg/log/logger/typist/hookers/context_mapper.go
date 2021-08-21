package hooks

import (
	"context"
	"reflect"

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
	_requestId := ctx.Value(keys.RequestId)
	if _requestId != nil && reflect.TypeOf(_requestId).Kind() == reflect.String {
		requestId := _requestId.(string)

		if len(requestId) != 0 {
			rs.Meta[keys.RequestId] = requestId
		}
	}

	return nil
}

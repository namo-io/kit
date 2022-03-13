package mctx

import (
	"context"

	"github.com/namo-io/kit/pkg/key"
)

func getStringFromContext(ctx context.Context, key key.Key) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}

	return v.(string)
}

func GetAuthorization(ctx context.Context) string {
	return getStringFromContext(ctx, key.Authorization)
}

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, key.Authorization, authorization)
}

func GetRequestId(ctx context.Context) string {
	return getStringFromContext(ctx, key.RequestId)
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, key.RequestId, requestId)
}

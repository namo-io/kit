package serverinterceptor

import (
	"context"
	"fmt"

	"github.com/namo-io/kit/pkg/ctxkey"
	"github.com/namo-io/kit/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var Default = []grpc.UnaryServerInterceptor{
	InjectContextByMetadata,
	ErrorHandling,
}

func InjectContextByMetadata(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata context is empty")
	}

	for _, key := range ctxkey.All {
		values := md.Get(key)
		if len(values) == 0 {
			continue
		}

		if len(values) == 1 {
			ctx = context.WithValue(ctx, key, values[0])
		} else {
			ctx = context.WithValue(ctx, key, values)
		}

	}

	return handler(ctx, req)
}

func ErrorHandling(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)

	resp, err := handler(ctx, req)
	if err != nil {
		gerr, ok := status.FromError(err)
		if ok {
			log.Error(gerr.Message())
			return resp, err
		} else {
			log.Error(err)
			return resp, status.New(codes.Internal, "An error occurred internally.").Err()
		}
	}

	return resp, err
}

package serverinterceptors

import (
	"context"
	"fmt"

	"github.com/namo-io/kit/pkg/keys"
	"github.com/namo-io/kit/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Default() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		InjectContextFromMD(),
		ErrorHandling(),
	}
}

func InjectContextFromMD() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("metadata context is empty")
		}

		for _, key := range keys.All() {
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
}

func ErrorHandling() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := log.WithContext(ctx)

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error(err)
		}

		_, ok := status.FromError(err)
		if ok {
			return resp, err
		}

		return resp, status.Error(codes.Internal, err.Error())
	}
}

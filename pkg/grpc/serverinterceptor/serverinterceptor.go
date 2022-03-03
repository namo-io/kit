package serverinterceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/namo-io/kit/pkg/ctxkey"
	"github.com/namo-io/kit/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var Default = []grpc.UnaryServerInterceptor{
	InjectContextByMetadata,
	InjectRequestID,
	ErrorHandling,
	Logging,
}

func InjectContextByMetadata(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
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
	}

	return handler(ctx, req)
}

func InjectRequestID(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestId := ctx.Value("x-request-id")
	if requestId == "" {
		requestId = ctx.Value(ctxkey.RequestId)
	}

	if requestId == "" {
		requestId = uuid.New().String()
	}

	ctx = context.WithValue(ctx, ctxkey.RequestId, requestId)
	return handler(ctx, req)
}

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)
	log.Debugf("request comming... method: '%v'", info.FullMethod)
	return handler(ctx, req)
}

func ErrorHandling(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)

	resp, err := handler(ctx, req)
	if err != nil {
		gerr, ok := status.FromError(err)
		if ok {
			log.Error(gerr.Message())
			return resp, gerr.Err()
		} else {
			log.Error(err)
			return resp, status.New(codes.Internal, "An error occurred internally.").Err()
		}
	}

	return resp, err
}

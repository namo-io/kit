package serverinterceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/namo-io/kit/pkg/log"
	"github.com/namo-io/kit/pkg/mctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var Default = []grpc.UnaryServerInterceptor{
	InjectContextAuthoriztion,
	InjectRequestID,
	Logging,
	Tracing,
	ErrorHandling,
}

func InjectContextAuthoriztion(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authorizationCode := md.Get(mctx.AuthorizationKey.String())
		if len(authorizationCode) == 1 {
			ctx = mctx.WithAuthorization(ctx, authorizationCode[0])
		}
	}

	return handler(ctx, req)
}

func InjectRequestID(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestId := ctx.Value("x-request-id")
	if requestId == "" {
		requestId = ctx.Value(mctx.RequestIdKey.String())
	}

	if requestId == "" {
		ctx = mctx.WithRequestId(ctx, uuid.New().String())
	}

	return handler(ctx, req)
}

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := log.WithContext(ctx)
	log.WithField("method", info.FullMethod).Debugf("grpc request comming...")
	return handler(ctx, req)
}

func Tracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	r, err := handler(ctx, req)
	log.Trace("Tracing end")

	return r, err
}

func ErrorHandling(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Trace("ErrorHandling")
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

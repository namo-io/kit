package serverinterceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/namo-io/kit/pkg/log"
	"github.com/namo-io/kit/pkg/mctx"
	"github.com/namo-io/kit/pkg/metric"
	"github.com/namo-io/kit/pkg/trace"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
	Metrics,
	ErrorHandling,
}

func InjectContextAuthoriztion(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headers := []string{"x-authorization", "authorization", "x-authorization-code", "authorization-code"}
		for _, header := range headers {
			val := md.Get(header)
			if len(val) == 1 {
				ctx = mctx.WithAuthorization(ctx, val[0])
			}
		}
	}

	return handler(ctx, req)
}

func InjectRequestID(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		headers := []string{"x-request-id", "request-id"}
		for _, header := range headers {
			val := md.Get(header)
			if len(val) == 1 {
				ctx = mctx.WithRequestId(ctx, val[0])
			}
		}

		if len(mctx.GetRequestId(ctx)) == 0 {
			ctx = mctx.WithRequestId(ctx, uuid.New().String())
		}
	}

	return handler(ctx, req)
}

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.WithContext(ctx).WithField("method", info.FullMethod).Debugf("grpc request comming...")
	return handler(ctx, req)
}

func Tracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := trace.Start(ctx, info.FullMethod)
	defer span.End()

	return handler(ctx, req)
}

func Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	resp, err := handler(ctx, req)
	latencySeconds := time.Now().Sub(t).Seconds()

	labels := prometheus.Labels{
		"method":         info.FullMethod,
		"status_code":    "",
		"exception_type": "",
	}

	if err != nil {
		gerr, ok := status.FromError(err)
		if ok {
			labels["status_code"] = fmt.Sprintf("%v", int(gerr.Code()))

			for _, detail := range gerr.Details() {
				switch d := detail.(type) {
				case *errdetails.ErrorInfo:
					labels["exception_type"] = d.GetDomain()
				}
			}
		}
	} else {
		labels["status_code"] = fmt.Sprintf("%v", int(codes.OK))
	}

	metric.GrpcRequestsTotal.With(labels).Add(1)
	metric.GrpcRequestDurationSeconds.With(labels).Observe(latencySeconds)

	return resp, err
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

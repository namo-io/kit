package clientinterceptor

import (
	"context"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var Default = []grpc.UnaryClientInterceptor{
	ErrorHandling,
}

func ErrorHandling(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)

	st, ok := status.FromError(err)
	if ok {
		for _, detail := range st.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				return fmt.Errorf("code: %s, domain: %s, description: %s", st.Code(), d.GetDomain(), st.Message())
			}
		}
	}

	return err
}

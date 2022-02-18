package err

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(s codes.Code, code string, description string) error {
	st, _ := status.New(s, description).WithDetails(&errdetails.ErrorInfo{
		Domain: code,
		Reason: description,
	})

	return st.Err()
}

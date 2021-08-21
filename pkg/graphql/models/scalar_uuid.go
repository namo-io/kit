package graphql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func MarshalUUIDScalar(id *uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.String()))
	})
}

func UnmarshalUUIDScalar(v interface{}) (*uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		uuid, err := uuid.Parse(v)
		return &uuid, err
	case *string:
		if v == nil {
			return nil, nil
		}

		uuid, err := uuid.Parse(*v)
		return &uuid, err
	default:
		return nil, fmt.Errorf("%T is not a uuid", v)
	}
}

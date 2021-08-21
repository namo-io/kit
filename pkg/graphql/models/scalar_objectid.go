package graphql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarshalObjectIDScalar(id *primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.Hex()))
	})
}

func UnmarshalObjectIDScalar(v interface{}) (*primitive.ObjectID, error) {
	switch v := v.(type) {
	case string:
		objectID, err := primitive.ObjectIDFromHex(v)
		return &objectID, err
	case *string:
		if v == nil {
			return nil, nil
		}

		objectID, err := primitive.ObjectIDFromHex(*v)
		return &objectID, err
	default:
		return nil, fmt.Errorf("%T is not a objectID", v)
	}
}

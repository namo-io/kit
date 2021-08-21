package objectidpb

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BSONObjectID primitive.ObjectID

func (u BSONObjectID) Marshal() ([]byte, error) {
	return u[:], nil
}

func (u BSONObjectID) MarshalTo(data []byte) (int, error) {
	return copy(data, (u)[:]), nil
}

func (u *BSONObjectID) Unmarshal(d []byte) error {
	copy((*u)[:], d)
	return nil
}

func (u *BSONObjectID) Size() int {
	return len(*u)
}

func (u *BSONObjectID) UnmarshalBSONValue(t bsontype.Type, d []byte) error {
	copy(u[:], d)
	return nil
}

func (u BSONObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.ObjectID, u[:], nil
}

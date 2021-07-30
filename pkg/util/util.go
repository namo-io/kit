package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/namo-io/kit/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}

	return hostname
}

func WaitSignal() os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer signal.Stop(sigs)

	return <-sigs
}

func StringToObjectId(str string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(str)
}

func MustStringToObjectId(str string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		log.Warnf("must string to objectID is error, msg: %v", err.Error())
		return primitive.NilObjectID
	}

	return objectId
}

func StringsToObjectIds(slice []string) ([]primitive.ObjectID, error) {
	objectIds := make([]primitive.ObjectID, 0)
	for _, item := range slice {
		objectId, err := primitive.ObjectIDFromHex(item)
		if err != nil {
			return nil, err
		}

		objectIds = append(objectIds, objectId)
	}

	return objectIds, nil
}

func ObjectIdToString(objectId primitive.ObjectID) string {
	return objectId.Hex()
}

func ObjectIdsToStrings(objectIds []primitive.ObjectID) []string {
	strs := make([]string, 0)
	for _, objectId := range objectIds {
		strs = append(strs, ObjectIdToString(objectId))
	}
	return strs
}

package eventbus

import (
	"fmt"
	"testing"
)

func TestEventBus(t *testing.T) {
	e := New()
	e.Subscribe("TEST", SubscribeTest)
	e.Publish("TEST", "QWD")
	e.Unsubscribe("TEST", SubscribeTest)
	e.Publish("TEST", "QWD")
}

func SubscribeTest(msg Msg) {
	fmt.Println(msg)
}

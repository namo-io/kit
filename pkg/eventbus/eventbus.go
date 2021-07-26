package eventbus

import (
	"reflect"
	"runtime"

	"github.com/namo-io/kit/pkg/log"
)

type Msg interface{}
type EventHandler func(msg Msg)

type EventBus interface {
	Publish(subject string, msg Msg)
	Subscribe(subject string, handler EventHandler)
	Unsubscribe(subject string, handler EventHandler)
}

type eventBus struct {
	handleChain map[string][]EventHandler
}

func New() *eventBus {
	return &eventBus{
		handleChain: make(map[string][]EventHandler),
	}
}

func (e *eventBus) Publish(subject string, msg Msg) {
	for _, handle := range e.handleChain[subject] {
		handle(msg)
	}
}

func (e *eventBus) Subscribe(subject string, handler EventHandler) {
	e.handleChain[subject] = append(e.handleChain[subject], handler)
}

func (e *eventBus) Unsubscribe(subject string, handler EventHandler) {
	for _index, _handler := range e.handleChain[subject] {
		if runtime.FuncForPC(reflect.ValueOf(_handler).Pointer()).Name() ==
			runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name() {
			e.handleChain[subject] = append(e.handleChain[subject][:_index], e.handleChain[subject][_index+1:]...)
			return
		}
	}

	log.Warnf("eventbus: Unsubscribe is not found subject, subject: %s", subject)
}

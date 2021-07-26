package logger

import (
	"context"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Trace(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Tracef(string, ...interface{})

	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
	WithContext(context.Context) Logger
}

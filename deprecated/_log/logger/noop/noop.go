package noop

import (
	"context"

	"github.com/namo-io/kit/pkg/log/logger"
)

type noop struct {
}

func New() *noop {
	return &noop{}
}

func (t *noop) WithField(key string, value interface{}) logger.Logger {
	return t
}

func (t *noop) WithFields(fields map[string]interface{}) logger.Logger {
	return t
}

func (t *noop) WithContext(ctx context.Context) logger.Logger {
	return t
}

func (t *noop) Panic(args ...interface{}) {
}

func (t *noop) Fatal(args ...interface{}) {
}

func (t *noop) Error(args ...interface{}) {
}

func (t *noop) Warn(args ...interface{}) {
}

func (t *noop) Info(args ...interface{}) {
}

func (t *noop) Debug(args ...interface{}) {
}

func (t *noop) Trace(args ...interface{}) {
}

func (t *noop) Panicf(format string, args ...interface{}) {
}

func (t *noop) Fatalf(format string, args ...interface{}) {
}

func (t *noop) Errorf(format string, args ...interface{}) {
}

func (t *noop) Warnf(format string, args ...interface{}) {
}

func (t *noop) Infof(format string, args ...interface{}) {
}

func (t *noop) Debugf(format string, args ...interface{}) {
}

func (t *noop) Tracef(format string, args ...interface{}) {
}

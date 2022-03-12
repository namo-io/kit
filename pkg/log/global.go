package log

import (
	"context"
	"fmt"
)

var (
	gLog = &log{
		callframeDepth: 1,
		verbose:        true,
		fields:         map[string]string{},
	}

	gRecorders = []*Recorder{
		NewDefaultRecorder(),
	}

	ErrRecorderParameterIsNil = fmt.Errorf("recorder parameter is nil")
	ErrRecorderIsAlreadyAdded = fmt.Errorf("the recorder is already added")
)

// SetVerbose get more information about logging, default: true
func SetVerbose(v bool) {
	gLog.verbose = v
}

// Addrecorder add recorder into global log
func AddRecorder(r *Recorder) error {
	if r == nil {
		return ErrRecorderParameterIsNil
	}

	if IsExistRecorder(r) {
		return ErrRecorderIsAlreadyAdded
	}

	WithField("recordersLength", fmt.Sprintf("%v", len(gRecorders))).Tracef("log recorder added")
	gRecorders = append(gRecorders, r)

	return nil
}

// IsExistRecorder check to exist Recorder
func IsExistRecorder(r *Recorder) bool {
	for _, recorder := range gRecorders {
		if recorder == r {
			return false
		}
	}

	return true
}

// SetField set field to global log
func SetField(k string, v string) {
	gLog = gLog.WithField(k, v).(*log)
}

// SetFields set field to global log
func SetFields(fields map[string]string) {
	gLog = gLog.WithFields(fields).(*log)
}

func Trace(a ...any) {
	gLog.Trace(a...)
}

func Tracef(format string, a ...any) {
	gLog.Tracef(format, a...)
}

func Debug(a ...any) {
	gLog.Debug(a...)
}

func Debugf(format string, a ...any) {
	gLog.Debugf(format, a...)
}

func Warn(a ...any) {
	gLog.Warn(a...)
}

func Warnf(format string, a ...any) {
	gLog.Warnf(format, a...)
}

func Info(a ...any) {
	gLog.Info(a...)
}

func Infof(format string, a ...any) {
	gLog.Infof(format, a...)
}

func Error(a ...any) {
	gLog.Error(a...)
}

func Errorf(format string, a ...any) {
	gLog.Errorf(format, a...)
}

func Fatal(a ...any) {
	gLog.Fatal(a...)
}

func Fatalf(format string, a ...any) {
	gLog.Fatalf(format, a...)
}

func WithField(key string, value string) Log {
	return gLog.WithField(key, value)
}

func WithFields(fields map[string]string) Log {
	return gLog.WithFields(fields)
}

func WithContext(ctx context.Context) Log {
	return gLog.WithContext(ctx)
}

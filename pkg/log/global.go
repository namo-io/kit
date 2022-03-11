package log

import (
	"context"
	"fmt"
)

var (
	glog = &log{
		callframeDepth: 1,
		verbose:        true,
		recorders: []*Recorder{
			NewDefaultRecorder(),
		},
	}

	ErrRecorderParameterIsNil = fmt.Errorf("recorder parameter is nil")
	ErrRecorderIsAlreadyAdded = fmt.Errorf("the recorder is already added")
)

// SetVerbose get more information about logging, default: true
func SetVerbose(v bool) {
	glog.verbose = v
}

// Addrecorder add recorder into global log
func AddRecorder(r *Recorder) error {
	if r == nil {
		return ErrRecorderParameterIsNil
	}

	if IsExistRecorder(r) {
		return ErrRecorderIsAlreadyAdded
	}

	Tracef("log recorder added, recorders length: %d", len(glog.recorders))
	glog.recorders = append(glog.recorders, r)

	return nil
}

// IsExistRecorder check to exist Recorder
func IsExistRecorder(r *Recorder) bool {
	for _, recorder := range glog.recorders {
		if recorder == r {
			return false
		}
	}

	return true
}

func Trace(a ...any) {
	glog.Trace(a...)
}

func Tracef(format string, a ...any) {
	glog.Tracef(format, a...)
}

func Debug(a ...any) {
	glog.Debug(a...)
}

func Debugf(format string, a ...any) {
	glog.Debugf(format, a...)
}

func Warn(a ...any) {
	glog.Warn(a...)
}

func Warnf(format string, a ...any) {
	glog.Warnf(format, a...)
}

func Info(a ...any) {
	glog.Info(a...)
}

func Infof(format string, a ...any) {
	glog.Infof(format, a...)
}

func Error(a ...any) {
	glog.Error(a...)
}

func Errorf(format string, a ...any) {
	glog.Errorf(format, a...)
}

func Fatal(a ...any) {
	glog.Fatal(a...)
}

func Fatalf(format string, a ...any) {
	glog.Fatalf(format, a...)
}

func WithField(key string, value string) Log {
	return glog.WithField(key, value)
}

func WithContext(ctx context.Context) Log {
	return glog.WithContext(ctx)
}

package log

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/namo-io/kit/pkg/mctx"
)

type Log interface {
	Trace(...any)
	Debug(...any)
	Info(...any)
	Warn(...any)
	Error(...any)
	Fatal(...any)

	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)

	WithField(string, string) Log
	WithFields(map[string]string) Log
	WithContext(context.Context) Log
}

type log struct {
	m sync.Mutex

	fields    map[string]string
	verbose   bool
	recorders []*Recorder

	callframeDepth int
}

func (l *log) Trace(a ...any) {
	l.record(TraceLevel, a...)
}

func (l *log) Tracef(format string, a ...any) {
	l.record(TraceLevel, fmt.Sprintf(format, a...))
}

func (l *log) Debug(a ...any) {
	l.record(DebugLevel, a...)
}

func (l *log) Debugf(format string, a ...any) {
	l.record(DebugLevel, fmt.Sprintf(format, a...))
}

func (l *log) Warn(a ...any) {
	l.record(WarnLevel, a...)
}

func (l *log) Warnf(format string, a ...any) {
	l.record(WarnLevel, fmt.Sprintf(format, a...))
}

func (l *log) Info(a ...any) {
	l.record(InfoLevel, a...)
}

func (l *log) Infof(format string, a ...any) {
	l.record(InfoLevel, fmt.Sprintf(format, a...))
}

func (l *log) Error(a ...any) {
	l.record(ErrorLevel, a...)
}

func (l *log) Errorf(format string, a ...any) {
	l.record(ErrorLevel, fmt.Sprintf(format, a...))
}

func (l *log) Fatal(a ...any) {
	l.record(FatalLevel, a...)
	os.Exit(1)
}

func (l *log) Fatalf(format string, a ...any) {
	l.record(FatalLevel, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func (l *log) WithField(key string, value string) Log {
	if len(l.fields[key]) != 0 {
		Warnf("log field is already exist, key='%v'", key)
	}

	copylog := l.deepcopy()
	copylog.fields[key] = value

	return copylog
}

func (l *log) WithFields(fields map[string]string) Log {
	var copylog Log
	copylog = l

	for k, v := range fields {
		copylog = copylog.WithField(k, v)
	}

	return copylog
}

func (l *log) WithContext(ctx context.Context) Log {
	copylog := l.deepcopy()
	for _, k := range mctx.Keys {
		value := ctx.Value(k)
		if value == nil {
			continue
		}

		if reflect.TypeOf(value).String() != "string" {
			WithField("key", k.String()).Warnf("context value is not string")
			continue
		}

		copylog = l.WithField(k.String(), value.(string)).(*log)
	}

	return copylog
}

func (l *log) deepcopy() *log {
	copyfields := make(map[string]string)
	for key, value := range l.fields {
		copyfields[key] = value
	}

	return &log{
		verbose:        l.verbose,
		fields:         copyfields,
		recorders:      l.recorders,
		callframeDepth: 0,
	}
}

func (l *log) record(level Level, a ...any) {
	// if verbose option, ignore specific level (trace, debug)
	if !l.verbose && level < WarnLevel {
		return
	}

	// collect log context information
	ts := time.Now()
	pcs := make([]uintptr, 10)
	depth := runtime.Callers(1, pcs)
	callFrames := runtime.CallersFrames(pcs[2+l.callframeDepth : depth])

	var frames []runtime.Frame
	for frame, again := callFrames.Next(); again; frame, again = callFrames.Next() {
		frames = append(frames, frame)
	}

	// syncing
	l.m.Lock()
	defer l.m.Unlock()

	// record
	for _, recorder := range l.recorders {
		err := recorder.record(ts, frames, level, fmt.Sprint(a...), l.fields)
		if err != nil {
			fmt.Println(err)
		}
	}
}

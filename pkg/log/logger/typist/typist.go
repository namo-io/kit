package typist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
	"unicode"

	deepcopy "github.com/barkimedes/go-deepcopy"
	"github.com/namo-io/kit/pkg/log/logger"
)

type typist struct {
	// ouptut is recording target ex) os.Stdout
	output io.Writer

	// level is output target (Fatal, Info, Debug ...)
	// ex) if level = InfoLevel, only logging (panic, fatal, info)
	level Level

	// formatter is output formatter
	formatter Formatter

	ctx                     context.Context
	callerIgnorePackageFile string

	hookers []Hooker

	meta map[string]interface{}
}

const maxCallerDeath = 20

/*
   New create typist instance
   '
   output:    os.Stdout
   level:     DebugLevel
   formatter  defaultFormatter
*/
func New(opts ...Options) *typist {
	t := &typist{
		output:                  os.Stdout,
		level:                   TraceLevel,
		ctx:                     context.Background(),
		callerIgnorePackageFile: "",
		meta:                    make(map[string]interface{}),
		formatter:               NewDefaultFormatter(),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func (t *typist) log(level Level, args ...interface{}) {
	pcs := make([]uintptr, maxCallerDeath)
	depth := runtime.Callers(1, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var _frames []runtime.Frame
	var _self_file string
	for frame, again := frames.Next(); again; frame, again = frames.Next() {
		if len(_self_file) == 0 {
			_self_file = frame.File
			continue
		}

		if _self_file != frame.File {
			if t.callerIgnorePackageFile != frame.File {
				_frames = append(_frames, frame)
			}
		}
	}

	rs := &Record{
		Time:    time.Now(),
		Level:   level,
		Context: t.ctx,
		Message: fmt.Sprint(args...),
		Meta:    t.meta,
		Frames:  _frames,
	}

	for _, hooker := range t.hookers {
		if err := hooker.Fire(t.ctx, t.level, rs); err != nil {
			t.output.Write([]byte(t.formatter.Format(&Record{
				Time:    time.Now(),
				Level:   ErrorLevel,
				Context: t.ctx,
				Message: fmt.Sprintf("typist: hooker is not working, name: %v, err: %v", hooker.Name(), err.Error()),
				Meta:    t.meta,
				Frames:  _frames,
			})))
		}
	}

	if level <= t.level {
		t.output.Write([]byte(t.formatter.Format(rs)))
	}
}

func (t *typist) clone() *typist {
	meta := deepcopy.MustAnything(t.meta).(map[string]interface{})

	return &typist{
		output:                  t.output,
		level:                   t.level,
		formatter:               t.formatter,
		callerIgnorePackageFile: t.callerIgnorePackageFile,
		hookers:                 t.hookers,
		ctx:                     t.ctx,
		meta:                    meta,
	}
}

func (t *typist) AddHooker(hooker Hooker) error {
	if hooker == nil {
		return errors.New("typist: hooker is nil pointer")
	}

	t.hookers = append(t.hookers, hooker)
	return nil
}

func (t *typist) AddField(key string, value interface{}) {
	t.meta[key] = value
}

func (t *typist) WithField(key string, value interface{}) logger.Logger {
	if key == "" {
		t.Warnf("WithField key parameter is empty, key: %v", key)
	}

	if len(key) > 0 && unicode.IsLower(rune(key[0])) {
		t.Warnf("WithField key parameter is not letter character, key: %v", key)
	}

	clone := t.clone()
	clone.meta[key] = value

	return clone
}

func (t *typist) WithFields(fields map[string]interface{}) logger.Logger {
	var clone logger.Logger
	for k, v := range fields {
		clone = t.WithField(k, v)
	}

	return clone
}

func (t *typist) WithContext(ctx context.Context) logger.Logger {
	clone := t.clone()
	clone.ctx = ctx

	return clone
}

func (t *typist) Panic(args ...interface{}) {
	t.log(PanicLevel, args...)
}

func (t *typist) Fatal(args ...interface{}) {
	t.log(FatalLevel, args...)
}

func (t *typist) Error(args ...interface{}) {
	t.log(ErrorLevel, args...)
}

func (t *typist) Warn(args ...interface{}) {
	t.log(WarnLevel, args...)
}

func (t *typist) Info(args ...interface{}) {
	t.log(InfoLevel, args...)
}

func (t *typist) Debug(args ...interface{}) {
	t.log(DebugLevel, args...)
}

func (t *typist) Trace(args ...interface{}) {
	t.log(TraceLevel, args...)
}

func (t *typist) Panicf(format string, args ...interface{}) {
	t.log(PanicLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Fatalf(format string, args ...interface{}) {
	t.log(FatalLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Errorf(format string, args ...interface{}) {
	t.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Warnf(format string, args ...interface{}) {
	t.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Infof(format string, args ...interface{}) {
	t.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Debugf(format string, args ...interface{}) {
	t.log(DebugLevel, fmt.Sprintf(format, args...))
}

func (t *typist) Tracef(format string, args ...interface{}) {
	t.log(TraceLevel, fmt.Sprintf(format, args...))
}

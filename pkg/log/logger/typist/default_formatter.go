package typist

import (
	"fmt"
	"path"
	"strings"

	"github.com/namo-io/kit/pkg/util/color"
)

type DefaultFormatter struct {
	IsUseColor bool
	StackPrint *StackPrint
}

type StackPrint struct {
	Level Level

	// Size is max: 20, ref: typist.MaxCallDepth
	Size int
}

const timeFormat = "2006-01-02 15:04:05.000 (MST)"

func levelColor(level Level) color.AnsiColor {
	_color := color.DefaultAnsiColor
	switch level {
	case PanicLevel:
		_color = color.RedAnsiColor
	case FatalLevel:
		_color = color.RedAnsiColor
	case ErrorLevel:
		_color = color.RedAnsiColor
	case WarnLevel:
		_color = color.YellowAnsiColor
	case InfoLevel:
		_color = color.CyanAnsiColor
	case DebugLevel:
		_color = color.BlueAnsiColor
	case TraceLevel:
		_color = color.GrayAnsiColor
	}

	return _color
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		IsUseColor: true,
		StackPrint: &StackPrint{
			Level: ErrorLevel,
			Size:  10,
		},
	}
}

func (f *DefaultFormatter) Format(rs *Record) (t string) {
	levelColor := levelColor(rs.Level)
	levelString := levelColor.Render(strings.ToUpper(rs.Level.String())[:4])
	t += fmt.Sprintf("%v: %v %v", rs.Time.Format(timeFormat), levelString, rs.Message)

	// func, file
	if len(rs.Frames) > 0 {
		funcname := "none"
		if funcnames := strings.Split(rs.Frames[0].Function, "."); len(funcnames) > 0 {
			funcname = funcnames[len(funcnames)-1]
		}
		filename := path.Base(rs.Frames[0].File)
		t += fmt.Sprintf(" %v=%v", levelColor.Render("Func"), fmt.Sprintf("%v(%d).%v", filename, rs.Frames[0].Line, funcname))
	}

	// print meta
	for k, v := range rs.Meta {
		t += fmt.Sprintf(" %v=%v", levelColor.Render(k), v)
	}

	// print caller stack
	if f.StackPrint != nil {
		if f.StackPrint.Level >= rs.Level {
			c := 0
			for _, frame := range rs.Frames {
				if c >= f.StackPrint.Size {
					break
				}

				if frame.Func != nil {
					t += fmt.Sprintf("%v    %v (0x%x)", NewLine, frame.Func.Name(), frame.Func.Entry())
					t += fmt.Sprintf("%v        %s:%d", NewLine, frame.File, frame.Line)
					c++
				}
			}
		}
	}

	t += NewLine

	return t
}

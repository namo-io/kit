package log

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type RecordOptions func(*Recorder)
type Record func(time.Time, []runtime.Frame, Level, string) error

type Recorder struct {
	record Record
}

func NewRecorder(r Record, opts ...RecordOptions) *Recorder {
	recorder := &Recorder{
		record: r,
	}

	for _, opt := range opts {
		opt(recorder)
	}

	return recorder
}

func NewDefaultRecorder() *Recorder {
	const timeFormat = "2006/01/02 15:04:05.000 (MST)"

	return NewRecorder(func(ts time.Time, frames []runtime.Frame, level Level, msg string) error {
		// funcname := "unknown"
		filename := "unknown"

		if len(frames) > 0 {
			if funcnames := strings.Split(frames[0].Function, "."); len(funcnames) > 0 {
				// funcname = funcnames[len(funcnames)-1]
			}
			filename = path.Base(frames[0].File)
		}

		_, err := fmt.Println(fmt.Sprintf("%v %v:%v [%v] %v",
			ts.Format(timeFormat),
			filename,
			frames[0].Line,
			NewColorByLogLevel(level).Render(strings.ToUpper(level.toString()[:4])),
			fmt.Sprint(msg),
		))
		return err
	})
}

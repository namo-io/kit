package typist

import (
	"context"
	"runtime"
	"time"
)

type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Context context.Context
	Meta    map[string]interface{}
	Frames  []runtime.Frame
}

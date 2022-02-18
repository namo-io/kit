package log

type Level int

func (l *Level) toString() string {
	switch *l {
	case 1:
		return "Trace"
	case 2:
		return "Debug"
	case 3:
		return "Warn"
	case 4:
		return "Info"
	case 5:
		return "Error"
	case 6:
		return "Fatal"
	default:
		return "Unknown"
	}
}

const (
	TraceLevel = 1 + iota
	DebugLevel
	WarnLevel
	InfoLevel
	ErrorLevel
	FatalLevel
)

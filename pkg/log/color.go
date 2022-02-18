package log

import "fmt"

// ANSI escape color code formatter
// support color tables docs: https://en.wikipedia.org/wiki/ANSI_escape_code
type AnsiColor int

const (
	DefaultAnsiColor = AnsiColor(-1)
	BlackAnsiColor   = AnsiColor(30)
	RedAnsiColor     = AnsiColor(31)
	YellowAnsiColor  = AnsiColor(33)
	BlueAnsiColor    = AnsiColor(34)
	CyanAnsiColor    = AnsiColor(36)
	GrayAnsiColor    = AnsiColor(90)
)

// Render is rendering text as color text
func (a AnsiColor) Render(t string) string {
	if a == DefaultAnsiColor {
		return t
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", int(a), t)
}

func NewColorByLogLevel(level Level) AnsiColor {
	color := DefaultAnsiColor
	switch level {
	case ErrorLevel:
		color = RedAnsiColor
	case WarnLevel:
		color = YellowAnsiColor
	case InfoLevel:
		color = CyanAnsiColor
	case DebugLevel:
		color = BlueAnsiColor
	case TraceLevel:
		color = GrayAnsiColor
	}

	return color
}

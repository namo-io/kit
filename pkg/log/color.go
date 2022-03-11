package log

import "fmt"

type Color int

const (
	DefaultColor   = Color(-1)
	BlackColor     = Color(232)
	RedColor       = Color(9)
	YellowColor    = Color(11)
	BlueColor      = Color(12)
	CyanColor      = Color(14)
	GrayColor      = Color(239)
	LightGrayColor = Color(245)
)

// Render is rendering text as color text
func (a Color) Render(t string) string {
	if a == DefaultColor {
		return t
	}

	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", int(a), t)
}

func newColorByLogLevel(level Level) Color {
	color := DefaultColor
	switch level {
	case ErrorLevel:
		color = RedColor
	case WarnLevel:
		color = YellowColor
	case InfoLevel:
		color = CyanColor
	case DebugLevel:
		color = BlueColor
	case TraceLevel:
		color = LightGrayColor
	}

	return color
}

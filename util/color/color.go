package color

import "fmt"

type ColorType int

const (
	red    ColorType = 31
	yellow ColorType = 33
	blue   ColorType = 36
	gray   ColorType = 37
)

func String(color ColorType, str string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, str)
}

func Yellow(str string) string {
	return String(yellow, str)
}

func Red(str string) string {
	return String(red, str)
}

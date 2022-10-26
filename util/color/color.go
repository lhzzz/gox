package color

import "fmt"

type ColorType int

const (
	Red    ColorType = 31
	Yellow ColorType = 33
	Blue   ColorType = 36
	Gray   ColorType = 37
)

func String(color ColorType, str string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, str)
}

func YollowStr(str string) string {
	return String(Yellow, str)
}

func RedStr(str string) string {
	return String(Red, str)
}

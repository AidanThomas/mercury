package log

import (
	"strings"
)

type Level int

const (
	Fatal Level = iota
	Error
	Info
	Debug
)

func (l *Level) Parse(s string) {
	switch strings.ToLower(s) {
	case "fatal":
		*l = Fatal
	case "error":
		*l = Error
	case "info":
		*l = Info
	case "debug":
		*l = Debug
	default:
		*l = Error
	}
}

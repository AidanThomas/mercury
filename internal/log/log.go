package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	timeStampFormat = "2006-01-02 15:04:05"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

func Errorf(format string, args ...interface{}) {
	format = addTimeStamp(format)
	stderr.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	var level Level
	level.Parse(os.Getenv("LOG_LEVEL"))
	if level >= Info {
		format = addTimeStamp(format)
		stdout.Printf(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	var level Level
	level.Parse(os.Getenv("LOG_LEVEL"))
	if level >= Debug {
		format = addTimeStamp(format)
		stdout.Printf(format, args...)
	}
}

func addTimeStamp(format string) string {
	return fmt.Sprintf("[%s] %s", time.Now().Format(timeStampFormat), format)
}

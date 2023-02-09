package seedcoin

import (
	"fmt"
	"sync"
)

var loggerOnceSyncPoint sync.Once

type Logger struct{}

var singletonLogger *Logger

const logTemplate string = "[SEEDCOIN]:"

func SharedLogger() *Logger {
	if singletonLogger == nil {
		loggerOnceSyncPoint.Do(
			func() {
				singletonLogger = &Logger{}
			},
		)
	}
	return singletonLogger
}

func (l *Logger) Log(format string, a ...any) {
	logFormat := logTemplate + " " + format
	logStr := fmt.Sprintf(logFormat, a)
	println(logStr)
}

func (l *Logger) LogMessage(message string) {
	logStr := logTemplate + " " + message
	println(logStr)
}

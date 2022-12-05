package seedcoin

import (
	"fmt"
	"sync"
)

var loggerOnceSyncPoint sync.Once

type Logger struct{}

var singletonLogger *Logger

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
	logTemplate := "[SEEDCOIN]: "
	logFormat := logTemplate + format
	logStr := fmt.Sprintf(logFormat, a)
	println(logStr)
}

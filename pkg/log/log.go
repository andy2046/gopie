// Package log implements a minimalistic Logger interface.
package log

import (
	"log"
	"os"
)

type (
	// Level defines Logging level.
	Level int

	// Logger is the Logging interface.
	Logger interface {
		Debug(v ...interface{})
		Info(v ...interface{})
		Warn(v ...interface{})
		Error(v ...interface{})
		Debugf(format string, v ...interface{})
		Infof(format string, v ...interface{})
		Warnf(format string, v ...interface{})
		Errorf(format string, v ...interface{})
	}

	defaultLogger struct {
		logger *log.Logger
		level  Level
	}
)

// Logging Levels.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// NewLogger returns a new Logger.
func NewLogger(level Level, prefix string) Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
		level:  level,
	}
}

func (dl *defaultLogger) log(level Level, v ...interface{}) {
	if dl.level <= level {
		dl.logger.Println(v...)
	}
}

func (dl *defaultLogger) logf(level Level, format string, v ...interface{}) {
	if dl.level <= level {
		dl.logger.Printf(format, v...)
	}
}

func (dl *defaultLogger) Debug(v ...interface{}) {
	dl.log(DEBUG, v...)
}

func (dl *defaultLogger) Info(v ...interface{}) {
	dl.log(INFO, v...)
}

func (dl *defaultLogger) Warn(v ...interface{}) {
	dl.log(WARN, v...)
}

func (dl *defaultLogger) Error(v ...interface{}) {
	dl.log(ERROR, v...)
}

func (dl *defaultLogger) Debugf(format string, v ...interface{}) {
	dl.logf(DEBUG, format, v...)
}

func (dl *defaultLogger) Infof(format string, v ...interface{}) {
	dl.logf(INFO, format, v...)
}

func (dl *defaultLogger) Warnf(format string, v ...interface{}) {
	dl.logf(WARN, format, v...)
}

func (dl *defaultLogger) Errorf(format string, v ...interface{}) {
	dl.logf(ERROR, format, v...)
}

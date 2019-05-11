// Package log implements a minimalistic Logger interface.
package log

import (
	"io"
	"log"
	"os"
	"sync"
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
		SetLevel(l Level)
	}

	// Config used to init Logger.
	Config struct {
		Level        Level
		Prefix       string
		DebugHandler io.Writer
		InfoHandler  io.Writer
		WarnHandler  io.Writer
		ErrorHandler io.Writer
	}

	// Option applies config to Logger Config.
	Option = func(*Config) error

	defaultLogger struct {
		loggerz map[Level]*log.Logger
		config  *Config
		sync.RWMutex
	}
)

// Logging Levels.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	// DefaultConfig is the default Logger Config.
	DefaultConfig = Config{
		Level:        INFO,
		Prefix:       "",
		DebugHandler: os.Stdout,
		InfoHandler:  os.Stdout,
		WarnHandler:  os.Stdout,
		ErrorHandler: os.Stderr,
	}
)

// NewLogger returns a new Logger.
func NewLogger(options ...Option) Logger {
	logConfig := DefaultConfig
	setOption(&logConfig, options...)

	return &defaultLogger{
		config:  &logConfig,
		loggerz: setLoggerz(&logConfig),
	}
}

// SetLevel sets the Logging level.
func (dl *defaultLogger) SetLevel(l Level) {
	dl.Lock()
	dl.config.Level = l
	dl.Unlock()
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

func (dl *defaultLogger) log(level Level, v ...interface{}) {
	dl.RLock()
	defer dl.RUnlock()
	if dl.config.Level <= level {
		dl.loggerz[level].Println(v...)
	}
}

func (dl *defaultLogger) logf(level Level, format string, v ...interface{}) {
	dl.RLock()
	defer dl.RUnlock()
	if dl.config.Level <= level {
		dl.loggerz[level].Printf(format, v...)
	}
}

func setOption(c *Config, options ...func(*Config) error) error {
	for _, opt := range options {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

func setLoggerz(logConfig *Config) map[Level]*log.Logger {
	loggerz := make(map[Level]*log.Logger)
	nonErrorLogger := log.New(os.Stdout, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(os.Stderr, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)

	loggerz[DEBUG] = nonErrorLogger
	loggerz[INFO] = nonErrorLogger
	loggerz[WARN] = nonErrorLogger
	loggerz[ERROR] = errorLogger

	if logConfig.DebugHandler != os.Stdout {
		loggerz[DEBUG] = log.New(logConfig.DebugHandler, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logConfig.InfoHandler != os.Stdout {
		loggerz[INFO] = log.New(logConfig.InfoHandler, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logConfig.WarnHandler != os.Stdout {
		loggerz[WARN] = log.New(logConfig.WarnHandler, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logConfig.ErrorHandler != os.Stderr {
		loggerz[ERROR] = log.New(logConfig.ErrorHandler, logConfig.Prefix, log.Ldate|log.Ltime|log.Lshortfile)
	}

	return loggerz
}

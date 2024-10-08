package logging

import (
	"log"

	"github.com/pkg/errors"
)

const (
	Silent int = 0
	Err    int = 1
	Warn   int = 2
	Info   int = 3
	Debug  int = 4
)

type SimpleLogging interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	// Err(format string, a ...interface{})
	Err(err error, message string)
}

type Logger struct {
	level int
}

func NewLogger(level int) *Logger {
	return &Logger{
		level: level,
	}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	if l.level >= Debug {
		log.Printf(format, a...)
	}
}

func (l *Logger) Info(format string, a ...interface{}) {
	if l.level >= Info {
		log.Printf(format, a...)
	}
}

func (l *Logger) Warn(format string, a ...interface{}) {
	if l.level >= Warn {
		log.Printf(format, a...)
	}
}

// TODO: use errors.Wrap(error, string)
func (l *Logger) Err(err error, message string) {
	if l.level >= Err {
		log.Printf("%v", errors.Wrap(err, message))
		// errors.Wrap(error, string)
	}
}

// func (l *Logger) Err(format string, a ...interface{}) {
// 	if l.level >= Err {
// 		log.Printf(format, a...)
// 	}
// }

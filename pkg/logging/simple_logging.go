package logging

import (
	"log"
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
	Err(format string, a ...interface{})
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

func (l *Logger) Err(format string, a ...interface{}) {
	if l.level >= Err {
		log.Printf(format, a...)
	}
}

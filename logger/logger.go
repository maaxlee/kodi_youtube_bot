package logger

import (
	"io"
	"log"
)

type Logger struct {
	log   *log.Logger
	debug bool
}

func GetLogger(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		log:   log.New(out, prefix, flag),
		debug: true,
	}
}

func (l *Logger) Printf(msg interface{}) {
	l.log.Print(msg)
}

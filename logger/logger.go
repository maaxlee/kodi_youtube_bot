package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	log   *log.Logger
	debug bool
}

func GetLogger(out io.Writer, prefix string, flag int) *Logger {
	debug := true
	debugStr, ok := os.LookupEnv("DEBUG")
	if !ok || debugStr != "1" {
		debug = false
	}
	return &Logger{
		log:   log.New(out, prefix, flag),
		debug: debug,
	}
}

func (l *Logger) Printf(msg interface{}) {
	l.log.Println(msg)
}

func (l *Logger) Debugp(msg interface{}) {
	if l.debug {
		l.log.Println(msg)

	}
}

package logger

import (
	"io"
	"log"
)

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Debugf(string, ...interface{})
	Warnf(string, ...interface{})
}

var levels = map[string]int{"off": 0, "error": 1, "info": 2, "warn": 3, "debug": 9, "all": 10}

var logger Logger
var level = "all"

func Infof(format string, args ...interface{}) {
	if levels[level] >= levels["info"] {
		log.Printf("[Info] "+format+"\n", args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if levels[level] >= levels["error"] {
		log.Printf("[Error] "+format+"\n", args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if levels[level] >= levels["debug"] {
		log.Printf("[Debug] "+format+"\n", args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if levels[level] >= levels["warn"] {
		log.Printf("[Warn] "+format+"\n", args...)
	}
}

func SetLevel(lev string) {
	level = lev
}

func SetOutput(out io.Writer) {
	log.SetOutput(out)
}

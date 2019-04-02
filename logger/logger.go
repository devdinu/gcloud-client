package logger

import (
	"log"

	"github.com/devdinu/gcloud-client/config"
)

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Debugf(string, ...interface{})
	Warnf(string, ...interface{})
}

var levels = map[string]int{"off": 0, "error": 1, "info": 2, "warn": 3, "debug": 9, "all": 10}

var logger Logger

func Infof(format string, args ...interface{}) {
	if levels[config.LogLevel()] >= levels["info"] {
		log.Printf("[Info] "+format+"\n", args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if levels[config.LogLevel()] >= levels["error"] {
		log.Printf("[Error] "+format+"\n", args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if levels[config.LogLevel()] >= levels["debug"] {
		log.Printf("[Debug] "+format+"\n", args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if levels[config.LogLevel()] >= levels["warn"] {
		log.Printf("[Warn] "+format+"\n", args...)
	}
}

func New() Logger {
	return logger
}

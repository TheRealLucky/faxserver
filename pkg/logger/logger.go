package logger

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}
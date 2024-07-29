package gormv2

import (
	"github.com/sirupsen/logrus"
)

type LoggerFunc func(...interface{})

func (f LoggerFunc) Print(args ...interface{}) { f(args...) }

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	//fmt.Printf(format + "\n", args...)
	logrus.Infof(format, args...)
}

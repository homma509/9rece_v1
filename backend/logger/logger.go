package logger

import (
	"log"
	"os"
)

// type logger interface {
// 	Infof(format string, args ...interface{})
// 	Warnf(format string, args ...interface{})
// 	Errorf(format string, args ...interface{})
// }

var l = &stdLogger{
	stdout: log.New(os.Stdout, "", 0),
	stderr: log.New(os.Stderr, "", 0),
}

type stdLogger struct {
	stderr *log.Logger
	stdout *log.Logger
}

func Infof(format string, args ...interface{}) {
	l.stdout.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	l.stderr.Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	l.stderr.Printf(format, args...)
}

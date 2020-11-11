package log

import (
	"log"
	"os"
	"strings"
)

var (
	// AppLogger Loggerのグローバルインスタンス
	AppLogger = newLogger()
)

// Logger Loggerインターフェス
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type appLogger struct {
	stderr *log.Logger
	stdout *log.Logger
}

func newLogger() Logger {
	return &appLogger{
		stdout: log.New(os.Stdout, "", 0),
		stderr: log.New(os.Stderr, "", 0),
	}
}

// Info Infoレベルのログ出力
func (l *appLogger) Info(args ...interface{}) {
	args = append([]interface{}{"Level", "INFO", "Message"}, args...)
	format := format(len(args) / 2)
	l.stdout.Printf(format, args...)
}

// Error Errorレベルのログ出力
func (l *appLogger) Error(args ...interface{}) {
	args = append([]interface{}{"Level", "ERROR", "Message"}, args...)
	format := format(len(args) / 2)
	l.stderr.Printf(format, args...)
}

func format(len int) string {
	var fields = make([]string, len, len)
	for i := range fields {
		fields[i] = "%s: %v"
	}
	return strings.Join(fields, "\t")
}

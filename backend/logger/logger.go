package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger Loggerの生成
func NewLogger() *zap.Logger {
	logger, err := newConfig().Build()
	if err != nil {
		log.Fatalf(
			"Level: %s, Msg: %s, Error: %v",
			"FATAL",
			"couldn't create logger config",
			err,
		)
	}
	return logger
}

func newConfig() *zap.Config {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	return &zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

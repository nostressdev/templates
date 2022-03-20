package logger

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(logLevel zapcore.Level) *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(logLevel))
	return logger
}
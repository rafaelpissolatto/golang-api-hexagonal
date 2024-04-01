package config

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"os"
)

// NewLogger Creates an instance of Uber Zap Logger with Elastic Common Schema (ECS) Logger Pattern
func NewLogger() *zap.SugaredLogger {
	encoder := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoder, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

// CloseLogger Sync flushes any buffered log entries.
func CloseLogger(log *zap.SugaredLogger) {
	_ = log.Sync()
}

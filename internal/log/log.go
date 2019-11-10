package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured *zap.Logger
func New() *zap.Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

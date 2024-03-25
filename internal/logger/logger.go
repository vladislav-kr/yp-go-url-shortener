package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

)

const (
	DebugLevel  = "debug"
	InfoLevel   = "info"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DPanicLevel = "dPanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
)

func MustLogger(level string) *zap.Logger {
	lvl := zap.InfoLevel

	switch level {
	case DebugLevel:
		lvl = zap.DebugLevel
	case InfoLevel:
		lvl = zap.InfoLevel
	case WarnLevel:
		lvl = zap.WarnLevel
	case ErrorLevel:
		lvl = zap.ErrorLevel
	case DPanicLevel:
		lvl = zap.DPanicLevel
	case PanicLevel:
		lvl = zap.PanicLevel
	case FatalLevel:
		lvl = zap.FatalLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(lvl),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

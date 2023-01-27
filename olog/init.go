package olog

import (
	"errors"
	"log"
	"runtime"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// InitZap logger
// Deprecated: this function is no longer needed
func InitLog() {
	var (
		logg *zap.Logger
		err  error
	)

	// if the log is already initialized, do nothing
	if logger != nil {
		return
	}

	cfg := zap.Config{
		Encoding:         "json",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			StacktraceKey: "stacktrace",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey: "caller",
			EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
				_, caller.File, caller.Line, _ = runtime.Caller(7)
				enc.AppendString(caller.FullPath())
			},
		},
	}

	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg.Development = true

	logg, err = cfg.Build()
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := logg.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Println(err)
		}
	}()

	// define logger
	logger = logg.Sugar()
}

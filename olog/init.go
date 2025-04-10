package olog

import (
	"log"
	"runtime"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// initLog logger internal library
func initLog() {
	var (
		logg *zap.Logger
		err  error
	)

	// if the log is already initialized, do nothing
	if Logger != nil {
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
	core := zapcore.NewTee(
		logg.Core(),
		otelzap.NewCore("OTOOLS-LOG", otelzap.WithLoggerProvider(global.GetLoggerProvider())),
	)
	logg = zap.New(core)
	if err != nil {
		log.Println(err)
	}

	// define logger
	Logger = logg.Sugar()
}

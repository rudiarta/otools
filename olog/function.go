package olog

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	IsExist := span.SpanContext().HasTraceID()
	if !IsExist {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// LogE error
//
//	func E(message interface{}) {
//		logger.Error(message)
//	}
func E(ctx context.Context, message interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Error(message)
}

// LogEf error with format
//
//	func Ef(format string, i ...interface{}) {
//		logger.Errorf(format, i...)
//	}
func Ef(ctx context.Context, format string, i ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Errorf(format, i...)
}

// LogI info
//
//	func I(message ...interface{}) {
//		logger.Info(message...)
//	}
func I(ctx context.Context, message ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Info(message...)
}

// LogIf info with format
//
//	func If(format string, i ...interface{}) {
//		logger.Infof(format, i...)
//	}
func If(ctx context.Context, format string, i ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Infof(format, i...)
}

// LogD info
//
//	func D(message ...interface{}) {
//		logger.Debug(message...)
//	}
func D(ctx context.Context, message ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Debug(message...)
}

// DF info with format
//
//	func DF(format string, i ...interface{}) {
//		logger.Debugf(format, i...)
//	}
func DF(ctx context.Context, format string, i ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Debugf(format, i...)
}

// W warn
//
//	func W(message ...interface{}) {
//		logger.Warn(message...)
//	}
func W(ctx context.Context, message ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Warn(message...)
}

// Wf warn with format
//
//	func Wf(format string, i ...interface{}) {
//		logger.Warnf(format, i...)
//	}
func Wf(ctx context.Context, format string, i ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Warnf(format, i...)
}

//	func Panic(i ...interface{}) {
//		logger.Panic(i)
//	}
func Panic(ctx context.Context, i ...interface{}) {
	if Logger == nil {
		initLog()
	}

	l := Logger.With(zap.String("trace-id", getTraceIDFromContext(ctx)), zap.Any("context", ctx))
	l.Panic(i)
}

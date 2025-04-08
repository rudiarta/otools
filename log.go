package otools

import (
	"context"
	"errors"
	llog "log"
	"os"
	"strings"
	"syscall"

	"github.com/rudiarta/otools/olog"
	"github.com/rudiarta/otools/otrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

var loggerProvider *log.LoggerProvider

func InitLog(host, serviceName, environment string) {
	ctx := context.Background()
	var err error

	// Create resource.
	res := otrace.NewResource(serviceName, environment)

	// Create a logger provider.
	// You can pass this instance directly when creating bridges.
	loggerProvider, err = newLoggerProvider(ctx, host, environment, res)
	if err != nil {
		panic(err)
	}

	// Register as global logger provider so that it can be accessed global.LoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(loggerProvider)
}

func newLoggerProvider(ctx context.Context, host, environment string, res *resource.Resource) (*log.LoggerProvider, error) {
	var exporter log.Exporter
	var err error
	conn := otrace.NewGrpcConn(context.Background(), host)

	switch {
	case strings.Contains(environment, "local"):
		f, err = os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			olog.E(context.Background(), err)
		}

		exporter, err = stdoutlog.New(
			stdoutlog.WithWriter(f),
		)
	case strings.Contains(environment, "test"):
	default:
		exporter, err = otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, err
		}
	}

	processor := log.NewBatchProcessor(exporter)
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	return provider, nil
}

func ShutDownLogProvider() error {
	if loggerProvider != nil {
		if err := loggerProvider.ForceFlush(context.Background()); err != nil {
			olog.DF(context.Background(), "Error flushing log provider: %v", err)
		}
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			olog.DF(context.Background(), "Error shutting down log provider: %v", err)
		}
		olog.D(context.Background(), "Shutting down & flushing log provider successfully")
	}

	err := olog.Logger.Desugar().Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		llog.Println(err)
	}

	return nil
}

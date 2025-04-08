package otools

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/rudiarta/otools/olog"
	"github.com/rudiarta/otools/otrace"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	otermetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var isInitMetric bool = false
var meter otermetric.Meter
var meterProvider *metric.MeterProvider

func ShutDownMeterProvider() error {
	if meterProvider != nil {
		if err := meterProvider.ForceFlush(context.Background()); err != nil {
			olog.DF(context.Background(), "Error flushing metric provider: %v", err)
		}
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			olog.DF(context.Background(), "Error shutting down metric provider: %v", err)
		}
		olog.D(context.Background(), "Shutting down & flushing metric provider successfully")
	}
	return nil
}

func InitMetrics(host, serviceName, environment string) error { // Not ready to use yet

	// Handle if env local or test metrics are not be exported
	switch {
	case strings.Contains(environment, "local"):
		var err error
		f, err = os.OpenFile("metric.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			olog.E(context.Background(), err)
		}

		exp, _ := stdoutmetric.New(
			stdoutmetric.WithWriter(f),
			// Use human-readable output.
			stdoutmetric.WithPrettyPrint(),
			// Do not print timestamps for the demo.
			stdoutmetric.WithoutTimestamps(),
		)

		reader := metric.NewPeriodicReader(exp,
			metric.WithInterval(30*time.Second),
			metric.WithTimeout(5*time.Second))
		// reader := metric.NewManualReader()

		meterProvider = metric.NewMeterProvider(
			metric.WithResource(otrace.NewResource(serviceName, environment)),
			metric.WithReader(reader),
		)

		otel.SetMeterProvider(meterProvider)

		meter = otel.Meter(
			"otools-metric-test",
			otermetric.WithInstrumentationVersion("v0.0.1"),
			otermetric.WithSchemaURL(semconv.SchemaURL),
		)

		err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			olog.E(context.Background(), err)
		}

		return nil
	case strings.Contains(environment, "test"):
		meter = noop.NewMeterProvider().Meter("otools-metric-test")
		return nil
	}

	conn := otrace.NewGrpcConn(context.Background(), host)

	// This reader is used as a stand-in for a reader that will actually export
	// data. See exporters in the go.opentelemetry.io/otel/exporters package
	// for more information.
	exp, err := otrace.NewExporterMetricGRPC(context.Background(), conn)
	if err != nil {
		println(err)
	}
	reader := metric.NewPeriodicReader(exp,
		metric.WithInterval(30*time.Second),
		metric.WithTimeout(5*time.Second))
	// reader := metric.NewManualReader()

	meterProvider = metric.NewMeterProvider(
		metric.WithResource(otrace.NewResource(serviceName, environment)),
		metric.WithReader(reader),
	)

	otel.SetMeterProvider(meterProvider)

	meter = otel.Meter(
		"otools-metric",
		otermetric.WithInstrumentationVersion("v0.0.1"),
		otermetric.WithSchemaURL(semconv.SchemaURL),
	)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		olog.E(context.Background(), err)
	}

	return nil
}

func HistogramMetric(metricName, metriCDescription, unitType string) otermetric.Float64Histogram {
	if !isInitMetric {
		meter = noop.NewMeterProvider().Meter("otools-metric-test")
	}

	if unitType == "" {
		unitType = "ms"
	}
	histogram, err := meter.Float64Histogram(metricName,
		otermetric.WithDescription(metriCDescription),
		otermetric.WithUnit(unitType))

	if err != nil {
		olog.E(context.Background(), err)
	}

	return histogram
}

func CounterMetric(metricName, metriCDescription, unitType string) otermetric.Int64Counter {
	if !isInitMetric {
		meter = noop.NewMeterProvider().Meter("otools-metric-test")
	}

	if unitType == "" {
		unitType = "1"
	}
	upDownCounter, err := meter.Int64Counter(metricName,
		otermetric.WithDescription(metriCDescription),
		otermetric.WithUnit(unitType))

	if err != nil {
		olog.E(context.Background(), err)
	}

	return upDownCounter
}

func UpDownCounterMetric(metricName, metriCDescription, unitType string) otermetric.Int64UpDownCounter {
	if !isInitMetric {
		meter = noop.NewMeterProvider().Meter("otools-metric-test")
	}

	if unitType == "" {
		unitType = "1"
	}
	upDownCounter, err := meter.Int64UpDownCounter(metricName,
		otermetric.WithDescription(metriCDescription),
		otermetric.WithUnit(unitType))

	if err != nil {
		olog.E(context.Background(), err)
	}

	return upDownCounter
}

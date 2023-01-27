// Package otrace ...
package otrace

import (
	"context"
	"io"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NewRequest struct {
	TraceID string
	SpanID  string
}

func NewGrpcConn(ctx context.Context, hostPort string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, hostPort,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithBlock(),
	)
	if err == nil {
		return conn
	}

	return nil
}

// newExporter returns a console exporter.
func NewExporterTraceGRPC(ctx context.Context, conn *grpc.ClientConn) (trace.SpanExporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func NewExporterTraceHttp(ctx context.Context, hostPort string) (trace.SpanExporter, error) {
	return otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(hostPort), otlptracehttp.WithInsecure())
}

func NewExporterTraceFile(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func NewExporterMetricGRPC(ctx context.Context, conn *grpc.ClientConn) (metric.Exporter, error) {
	return otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
}

func NewExporterMetricHttp(ctx context.Context, hostPort string) (metric.Exporter, error) {
	return otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(hostPort), otlpmetrichttp.WithInsecure())
}

func NewResource(serviceName, env string) *resource.Resource {
	resources, _ := resource.New(context.Background(),
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
	)
	r, _ := resource.Merge(
		resource.Default(),
		resources,
	)
	r, _ = resource.Merge(r, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.DeploymentEnvironmentKey.String(env),
		semconv.TelemetrySDKLanguageGo,
	))
	return r
}

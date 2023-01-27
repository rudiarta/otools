package otools

import (
	"context"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rudiarta/otools/olog"
	"github.com/rudiarta/otools/otrace"
	"github.com/rudiarta/otools/outils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var lpTracehost, lpTraceServiceName, lpTraceEnvironment string
var isInitTrace bool = false
var tp *tracesdk.TracerProvider
var f *os.File
var toolName string = "otools"

func constructNewSpanContext(ctx context.Context, request trace.SpanContext) (spanContext trace.SpanContext, err error) {
	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(request.TraceID().String())
	if err != nil {
		olog.Ef(ctx, "error: %s", err.Error())
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(request.SpanID().String())
	if err != nil {
		olog.Ef(ctx, "error: %s", err.Error())
		return spanContext, err
	}
	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = false
	spanContext = trace.NewSpanContext(spanContextConfig)
	return spanContext, nil
}

func constructNewSpanContextWithString(ctx context.Context, request otrace.NewRequest) (spanContext trace.SpanContext, err error) {
	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(request.TraceID)
	if err != nil {
		olog.Ef(ctx, "error: %s", err.Error())
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(request.SpanID)
	if err != nil {
		olog.Ef(ctx, "error: %s", err.Error())
		return spanContext, err
	}
	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = false
	spanContext = trace.NewSpanContext(spanContextConfig)
	return spanContext, nil
}

// host for otel-collecter GRPC Ex: "localhost:30080"
// serviceName Ex: "name_service"
// environment Ex: "DEV"
func getGrpcOtelTraceProvider(host, serviceName, environment string) (*tracesdk.TracerProvider, error) {
	ctx := context.Background()
	var exp tracesdk.SpanExporter

	conn := otrace.NewGrpcConn(ctx, host)

	switch {
	case strings.Contains(lpTraceEnvironment, "local"):
		// Write telemetry data to a file.
		var err error
		f, err = os.OpenFile("traces.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			olog.E(ctx, err)
		}
		exp, _ = otrace.NewExporterTraceFile(f)
	case strings.Contains(lpTraceEnvironment, "test"):
		exp = tracetest.NewNoopExporter()
	default:
		exp, _ = otrace.NewExporterTraceGRPC(ctx, conn)
	}

	res := otrace.NewResource(serviceName, environment)

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := tracesdk.NewBatchSpanProcessor(exp)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	)

	return tracerProvider, nil
}

// Deprecated: This function will be removed in a future release
// SetError func
// this function is still not implemented
func SetErrorFromContext(ctx context.Context, err error) {
	// span := trace.SpanFromContext(ctx)
	// if span == nil || err == nil {
	// 	return
	// }

	// span.SetAttributes(semconv.ExceptionTypeKey.String("ERROR"))
	// span.SetAttributes(semconv.ExceptionMessageKey.String(err.Error()))
	// span.SetAttributes(semconv.ExceptionStacktraceKey.String(string(debug.Stack())))
}

// GetTraceID func
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	IsExist := span.SpanContext().HasTraceID()
	if !IsExist {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// GetSpanID func
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	IsExist := span.SpanContext().HasSpanID()
	if !IsExist {
		return ""
	}

	return span.SpanContext().SpanID().String()
}

type tracerImpl struct {
	ctx  context.Context
	span trace.Span
	tags map[string]interface{}
}

type Tracer interface {
	Context() context.Context
	Tags() map[string]interface{}
	SetError(err error)
	Finish(additionalTags ...map[string]interface{})
}

// host for otel-collecter GRPC Ex: "localhost:30080"
// serviceName Ex: "name_service"
// environment Ex: "DEV"
func InitTracer(host, serviceName, environment string) {
	lpTracehost = host
	lpTraceServiceName = serviceName
	lpTraceEnvironment = environment
	tp, _ = getGrpcOtelTraceProvider(lpTracehost, lpTraceServiceName, lpTraceEnvironment)

	if !strings.Contains(lpTraceEnvironment, "test") {
		otel.SetTracerProvider(tp)
	}
	isInitTrace = true
}

func ShutDownTraceProvider() error {
	if tp != nil {
		if err := tp.Shutdown(context.Background()); err != nil {
			olog.DF(context.Background(), "Error shutting down tracer provider: %v", err)
		}
		olog.D(context.Background(), "Shutting down tracer provider successfully")
	}
	return nil
}

func StartTrace(ctx context.Context, operationName string) Tracer {

	if lpTracehost == "" {
		olog.I(ctx, "InitTracer first")
	}

	if tp == nil {
		olog.I(ctx, "InitTracer first")
	}

	var tr trace.Tracer
	switch {
	case strings.Contains(lpTraceEnvironment, "test") || !isInitTrace:
		tr = trace.NewNoopTracerProvider().Tracer(toolName)
	default:
		tr = tp.Tracer(toolName)
	}

	ctx, span := tr.Start(ctx, operationName)

	return &tracerImpl{
		ctx:  ctx,
		span: span,
	}
}

// Start tracer that will use new ctx from context.Background
func StartTracerWithContextBackground(parentCtx context.Context, operationName string) Tracer {
	span := trace.SpanFromContext(parentCtx)
	spanContext, err := constructNewSpanContext(parentCtx, span.SpanContext())
	if err != nil {
		olog.If(parentCtx, "ERROR: %s ", err.Error())
	}
	if ok := spanContext.IsValid(); !ok {
		olog.If(parentCtx, "IS VALID? : %v ", ok)
	}

	requestContext := context.Background()
	requestContext = trace.ContextWithSpanContext(requestContext, spanContext)

	return StartTrace(requestContext, operationName)
}

// Deprecated: This function will be removed in a future release
// Start tracer that will use new ctx from context.Background & using specific trace & span ID
func StartTracerWithTraceIDAndSpanIDContextBackground(ctx context.Context, operationName, TraceID, SpanID string) Tracer {
	var span trace.Span
	spanContext, err := constructNewSpanContextWithString(ctx, otrace.NewRequest{
		TraceID: TraceID,
		SpanID:  SpanID,
	})
	if err != nil {
		olog.If(ctx, "ERROR: %s ", err.Error())
	}
	if ok := spanContext.IsValid(); !ok {
		olog.If(ctx, "IS VALID? : %v ", ok)
	}

	requestContext := trace.ContextWithSpanContext(ctx, spanContext)

	if lpTracehost == "" {
		olog.I(ctx, "InitTracer first")
	}

	if tp == nil {
		olog.I(ctx, "InitTracer first")
	}

	var tr trace.Tracer
	switch {
	case strings.Contains(lpTraceEnvironment, "test") || !isInitTrace:
		tr = noop.NewTracerProvider().Tracer(toolName)
	default:
		tr = tp.Tracer(toolName)
	}

	ctx, span = tr.Start(requestContext, operationName, trace.WithSpanKind(trace.SpanKindServer))

	return &tracerImpl{
		ctx:  ctx,
		span: span,
	}
}

// Start tracer that will use specific trace & span ID
func StartTracerWithTraceIDAndSpanID(ctx context.Context, operationName, TraceID, SpanID string) Tracer {
	var span trace.Span
	spanContext, err := constructNewSpanContextWithString(ctx, otrace.NewRequest{
		TraceID: TraceID,
		SpanID:  SpanID,
	})
	if err != nil {
		olog.If(ctx, "ERROR: %s ", err.Error())
	}
	if ok := spanContext.IsValid(); !ok {
		olog.If(ctx, "IS VALID? : %v ", ok)
	}

	requestContext := trace.ContextWithSpanContext(ctx, spanContext)

	if lpTracehost == "" {
		olog.I(ctx, "InitTracer first")
	}

	if tp == nil {
		olog.I(ctx, "InitTracer first")
	}

	var tr trace.Tracer
	switch {
	case strings.Contains(lpTraceEnvironment, "test") || !isInitTrace:
		tr = noop.NewTracerProvider().Tracer(toolName)
	default:
		tr = tp.Tracer(toolName)
	}

	ctx, span = tr.Start(requestContext, operationName, trace.WithSpanKind(trace.SpanKindServer))

	return &tracerImpl{
		ctx:  ctx,
		span: span,
	}
}

// Context get active context
func (t *tracerImpl) Context() context.Context {
	return t.ctx
}

// Tags create tags in tracer span
func (t *tracerImpl) Tags() map[string]interface{} {
	t.tags = make(map[string]interface{})
	return t.tags
}

// SetError set error in span
// this function is still not implemented
func (t *tracerImpl) SetError(err error) {
	// SetErrorFromContext(t.ctx, err)
}

func (t *tracerImpl) Finish(tags ...map[string]interface{}) {
	defer func() {
		t.span.End(trace.WithStackTrace(true))
	}()

	// Debug trace set to default attibute
	t.span.SetAttributes(
		semconv.ExceptionTypeKey.String("DEBUG"),
		semconv.ExceptionMessageKey.String("Stack Trace Information"),
		semconv.ExceptionStacktraceKey.String(string(debug.Stack())),
	)

	if tags != nil && t.tags == nil {
		t.tags = make(map[string]interface{})
	}

	for _, tag := range tags {
		for k, v := range tag {
			t.tags[k] = v
		}
	}

	for k, v := range t.tags {
		var myKey = attribute.Key(k)
		t.span.SetAttributes(myKey.String(outils.ToString(v)))
	}
}

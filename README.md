# How To Use otools

## Installing SDK

Requirement: 
- GO >= 1.24
- go.opentelemetry.io/otel v1.35.0

```bash
go get github.com/rudiarta/otools@v0.0.3
```

## Init Metrics

```go
    import (
        "go.opentelemetry.io/otel/attribute"
        "github.com/rudiarta/otools"
    )


    // Add this code to your project first running
    // host for otel-collecter GRPC Ex: "localhost:30080"
    // serviceName Ex: "name_service"
    // environment Ex: "DEV"
    // Or
    // Environment Ex: set with prefix "local" 
    // if your local there is no otel-collector daemon running
    otools.InitMetrics(host, serviceName, environment)

    // Histogram
    hg := otools.HistogramMetric("company.CreatePickup", "create pickup histogram", "ms")
    hg.Record(ctx, value, attribute.String("metricType", "error"), attribute.String("url", "v1/ship/company/notify/{shipID}"))


    // Counter
    attrs := []attribute.KeyValue{attribute.String("metricType", "success"), attribute.String("value", "true")}
    ct := otools.CounterMetric("total.running.counter.goroutine", "untuk tau jumlah go routine", "1")
    ct.Add(inCtx, 1, attrs...)

    // Don't forget to execute this in graceful shutdown mode
    otools.ShutDownMeterProvider()
```

## Init Tracer

```go
    import "github.com/rudiarta/otools"

    // Add this code to your project first running
    // host for otel-collecter GRPC Ex: "localhost:30080"
    // serviceName Ex: "name_service"
    // environment Ex: "DEV"
    // Or
    // Environment Ex: set with prefix "local" 
    // if your local there is no otel-collector daemon running
    otools.InitTracer(host, serviceName, environment)

    // put this code in top of your function
    tt := otools.StartTrace(ctx, "operationName")
    ctx = tt.Context()
    defer tt.Finish(map[string]interface{}{
        "req": "...",
        "resp": "...",
        "...": any,
        })

    // passing context without get response `context canceled`
    go func() {
        ti := otools.StartTracerWithContextBackground(ctx, "operationName inner goroutine")
        inCtx := ti.Context() // inCtx not inherit deadline from ctx anymore
        defer ti.Finish(map[string]interface{}{
        "req": "...",
        "resp": "...",
        "...": any,
        })
    }()

    // Don't forget to execute this in graceful shutdown mode
    otools.ShutDownTraceProvider()
```

## Init Log
* New Update: log integration with log provider otelzap

```go
    import "github.com/rudiarta/otools/olog"
    
    // Add this code to your project first running
    // host for otel-collecter GRPC Ex: "localhost:30080"
    // serviceName Ex: "name_service"
    // environment Ex: "DEV"
    // Or
    // Environment Ex: set with prefix "local" 
    // if your local there is no otel-collector daemon running
    otools.InitLog(host, serviceName, environment)

    // Use context from ctx = tt.Context()
    // Generate by Tracer
    olog.E(ctx, any)
    olog.I(ctx, any)

     // Don't forget to execute this in graceful shutdown mode
    otools.ShutDownLogProvider()
```

## Http Request

```go
import "github.com/rudiarta/otools/outils"

res, err := outils.HTTPRequestJSON(&outils.HttpRequestParams{
		Context: context,
		Method:  "GET",
		URL:     "https://google.com",
		Header:  map[string]string{},
		Timeout: 30 * time.Second,
})
```
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Setup reads standard OTEL_* settings. The Collector resolves application
// identity from the pod UID supplied by Kubernetes, for every signal.
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	stdout := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(stdout))
	shutdown = func(context.Context) error { return nil }
	if os.Getenv("OTEL_SDK_DISABLED") == "true" || !hasEndpoint() {
		return shutdown, nil
	}

	var closers []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var errs []error
		for _, close := range closers {
			errs = append(errs, close(ctx))
		}
		return errors.Join(errs...)
	}
	defer func() {
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = errors.Join(err, shutdown(cleanupCtx))
		}
	}()

	res, err := resource.New(ctx, resource.WithTelemetrySDK(), resource.WithFromEnv())
	if err != nil {
		return shutdown, err
	}
	traces, err := otlptracehttp.New(ctx)
	if err != nil {
		return shutdown, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res), sdktrace.WithBatcher(traces))
	closers = append(closers, tp.Shutdown)

	metrics, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return shutdown, err
	}
	interval := 60 * time.Second
	if raw := os.Getenv("OTEL_METRIC_EXPORT_INTERVAL"); raw != "" {
		ms, parseErr := strconv.Atoi(raw)
		if parseErr != nil || ms <= 0 {
			return shutdown, fmt.Errorf("invalid OTEL_METRIC_EXPORT_INTERVAL: %q", raw)
		}
		interval = time.Duration(ms) * time.Millisecond
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithResource(res), sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(metrics, sdkmetric.WithInterval(interval)),
	))
	closers = append(closers, mp.Shutdown)

	logs, err := otlploghttp.New(ctx)
	if err != nil {
		return shutdown, err
	}
	lp := sdklog.NewLoggerProvider(sdklog.WithResource(res), sdklog.WithProcessor(sdklog.NewBatchProcessor(logs)))
	closers = append(closers, lp.Shutdown)

	// SDK/export errors stay on stderr to avoid an OTLP logging feedback loop.
	diagnostics := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		diagnostics.Error("OpenTelemetry export failed", "error", err)
	}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	global.SetLoggerProvider(lp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	slog.SetDefault(slog.New(fanoutHandler{stdout, otelslog.NewHandler("application", otelslog.WithLoggerProvider(lp))}))
	if err := runtime.Start(runtime.WithMeterProvider(mp)); err != nil {
		return shutdown, err
	}
	// Clients using the default transport, including OAuth clients, get spans.
	http.DefaultTransport = otelhttp.NewTransport(http.DefaultTransport)
	return shutdown, nil
}

func hasEndpoint() bool {
	for _, key := range []string{"OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "OTEL_EXPORTER_OTLP_LOGS_ENDPOINT"} {
		if os.Getenv(key) != "" {
			return true
		}
	}
	return false
}

type fanoutHandler []slog.Handler

func (h fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h fanoutHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range h {
		if handler.Enabled(ctx, record.Level) {
			errs = append(errs, handler.Handle(ctx, record.Clone()))
		}
	}
	return errors.Join(errs...)
}

func (h fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	result := make(fanoutHandler, len(h))
	for i, handler := range h {
		result[i] = handler.WithAttrs(attrs)
	}
	return result
}

func (h fanoutHandler) WithGroup(name string) slog.Handler {
	result := make(fanoutHandler, len(h))
	for i, handler := range h {
		result[i] = handler.WithGroup(name)
	}
	return result
}

package telemetry

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	collectorlog "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	collectormetric "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collectortrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
)

func TestHTTPExportsCorrelatedSignals(t *testing.T) {
	var mu sync.Mutex
	requests := make(map[string][][]byte)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		requests[r.URL.Path] = append(requests[r.URL.Path], body)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "k8s.pod.uid=test-pod,k8s.container.name=app")
	t.Setenv("OTEL_SDK_DISABLED", "false")
	t.Setenv("OTEL_SEMCONV_STABILITY_OPT_IN", "http")
	oldTracer, oldMeter, oldLogger := otel.GetTracerProvider(), otel.GetMeterProvider(), global.GetLoggerProvider()
	oldSlog, oldTransport := slog.Default(), http.DefaultTransport
	defer func() {
		otel.SetTracerProvider(oldTracer)
		otel.SetMeterProvider(oldMeter)
		global.SetLoggerProvider(oldLogger)
		slog.SetDefault(oldSlog)
		http.DefaultTransport = oldTransport
	}()
	shutdown, err := Setup(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			t.Error(err)
		}
	}()
	router := gin.New()
	router.Use(otelgin.Middleware(""), HTTPLogs())
	router.GET("/items/:id", func(c *gin.Context) { c.String(200, "ok") })
	req := httptest.NewRequest("GET", "/items/123?private=value", nil)
	req.Header.Set("traceparent", "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if recorder.Code != 200 {
		t.Fatalf("HTTP status: %d", recorder.Code)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	shutdown = func(context.Context) error { return nil }
	mu.Lock()
	defer mu.Unlock()
	var traceID, logTraceID string
	var hasDuration, hasPodUID bool
	for _, body := range requests["/v1/traces"] {
		var exported collectortrace.ExportTraceServiceRequest
		if err := proto.Unmarshal(body, &exported); err != nil {
			t.Fatal(err)
		}
		for _, resource := range exported.ResourceSpans {
			for _, attr := range resource.Resource.Attributes {
				if attr.Key == "k8s.pod.uid" && attr.Value.GetStringValue() == "test-pod" {
					hasPodUID = true
				}
			}
			for _, scope := range resource.ScopeSpans {
				for _, span := range scope.Spans {
					if span.Name == "GET /items/:id" {
						traceID = string(span.TraceId)
					}
				}
			}
		}
	}
	for _, body := range requests["/v1/logs"] {
		var exported collectorlog.ExportLogsServiceRequest
		if err := proto.Unmarshal(body, &exported); err != nil {
			t.Fatal(err)
		}
		for _, resource := range exported.ResourceLogs {
			for _, scope := range resource.ScopeLogs {
				for _, record := range scope.LogRecords {
					if record.Body.GetStringValue() == "HTTP request completed" {
						logTraceID = string(record.TraceId)
					}
				}
			}
		}
	}
	for _, body := range requests["/v1/metrics"] {
		var exported collectormetric.ExportMetricsServiceRequest
		if err := proto.Unmarshal(body, &exported); err != nil {
			t.Fatal(err)
		}
		for _, resource := range exported.ResourceMetrics {
			for _, scope := range resource.ScopeMetrics {
				for _, metric := range scope.Metrics {
					if metric.Name == "http.server.request.duration" {
						hasDuration = true
					}
					for _, point := range metric.GetHistogram().GetDataPoints() {
						for _, attr := range point.Attributes {
							if attr.Key == "trace_id" || attr.Value.GetStringValue() == "/items/123" {
								t.Fatalf("unbounded metric attribute: %v", attr)
							}
						}
					}
				}
			}
		}
	}
	if traceID == "" || logTraceID != traceID {
		t.Fatal("trace and log must preserve the incoming trace context")
	}
	if !hasDuration || !hasPodUID {
		t.Fatalf("missing HTTP metric or Kubernetes association: duration=%v pod=%v", hasDuration, hasPodUID)
	}
}

func TestDisabledWithoutCollector(t *testing.T) {
	for _, key := range []string{"OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "OTEL_EXPORTER_OTLP_LOGS_ENDPOINT"} {
		t.Setenv(key, "")
	}
	oldLogger := slog.Default()
	defer slog.SetDefault(oldLogger)
	shutdown, err := Setup(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

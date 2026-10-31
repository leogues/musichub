package telemetry

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func HTTPLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Request.URL.Path == "/healthz" {
			return
		}
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		span := trace.SpanContextFromContext(c.Request.Context())
		slog.InfoContext(c.Request.Context(), "HTTP request completed",
			"http.request.method", c.Request.Method, "http.route", route,
			"http.response.status_code", c.Writer.Status(), "duration.seconds", time.Since(start).Seconds(),
			"trace_id", span.TraceID().String(), "span_id", span.SpanID().String())
	}
}

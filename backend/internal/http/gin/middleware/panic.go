package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"go.opentelemetry.io/otel/trace"
)

func (m *Middleware) ReportPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				trace.SpanFromContext(c.Request.Context()).RecordError(fmt.Errorf("%v", err))
				slog.ErrorContext(c.Request.Context(), "HTTP handler panic", "error", err)
				musichub.ReportPanic(err)
			}
		}()
		c.Next()
	}
}

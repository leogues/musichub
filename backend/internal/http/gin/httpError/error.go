package httpError

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"go.opentelemetry.io/otel/trace"
)

func NewErrorHandler() func(c *gin.Context, err error) {
	return Error
}

var codes = map[string]int{
	musichub.ECONFLICT:       http.StatusConflict,
	musichub.EINVALID:        http.StatusBadRequest,
	musichub.ENOTFOUND:       http.StatusNotFound,
	musichub.ENOTIMPLEMENTED: http.StatusNotImplemented,
	musichub.EUNAUTHORIZED:   http.StatusUnauthorized,
	musichub.EINTERNAL:       http.StatusInternalServerError,
}

func Error(c *gin.Context, err error) {
	code, message := musichub.ErrorCode(err), musichub.ErrorMessage(err)
	if code == musichub.EINTERNAL {
		trace.SpanFromContext(c.Request.Context()).RecordError(err)
		musichub.ReportError(c.Request.Context(), err, c)
		LogError(c, err)
	}

	c.JSON(ErrorStatusCode(code), gin.H{"error": message})
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func LogError(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), "HTTP request failed", "http.request.method", c.Request.Method, "http.route", c.FullPath(), "error", err)
}

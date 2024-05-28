package httpError

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
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
	log.Printf("[http] error: %s %s: %s ", c.Request.Method, c.Request.URL.Path, err)
}

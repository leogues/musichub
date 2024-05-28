package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
)

func (m *Middleware) ReportPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				musichub.ReportPanic(err)
			}
		}()
		c.Next()
	}
}

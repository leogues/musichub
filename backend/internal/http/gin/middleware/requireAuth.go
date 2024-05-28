package middleware

import (
	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
)

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if userID := musichub.UserIDFromContext(ctx); userID != 0 {
			c.Next()
			return
		}

		m.Error(c, musichub.Errorf(musichub.EUNAUTHORIZED, "Unauthorized access"))
		c.Abort()
	}
}

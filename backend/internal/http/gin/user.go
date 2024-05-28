package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
)

func (s *GinServer) registerUserRoutes(r *gin.RouterGroup) {
	r.GET("/me", s.getMe)
	r.POST("/auth/logout", s.logout)
}

func (s *GinServer) logout(c *gin.Context) {
	if err := s.sessionManager.SetSession(c, session.NewSession()); err != nil {
		s.Error(c, fmt.Errorf("cannot set session cookie: %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (s *GinServer) getMe(c *gin.Context) {
	ctx := c.Request.Context()
	user := musichub.UserFromContext(ctx)

	c.JSON(200, user)
}

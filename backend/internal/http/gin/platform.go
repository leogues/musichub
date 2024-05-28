package gin

import "github.com/gin-gonic/gin"

func (s *GinServer) registerPlatformRoutes(r *gin.RouterGroup) {
	r.GET("/platforms/list", s.getPlaforms)
}

func (s *GinServer) getPlaforms(c *gin.Context) {
	platforms := s.ProviderService.Platforms()

	c.JSON(200, platforms)
}

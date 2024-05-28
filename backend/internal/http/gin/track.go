package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

func (s *GinServer) withTrackService(f provider.TrackProviderFactoryInterface, handleFunc func(*gin.Context, musichub.TrackService)) gin.HandlerFunc {
	return func(c *gin.Context) {
		providerName := c.Param("provider")

		token := musichub.ProviderTokenFromContext(c.Request.Context())

		service, err := f.CreateTrackProvider(providerName, token)

		if err != nil {
			s.Error(c, err)
			return
		}

		handleFunc(c, service)
	}
}

func (s *GinServer) registerTrackRoutes(r *gin.RouterGroup, f provider.TrackProviderFactoryInterface) {
	r.GET("/me/platforms/:provider/track", s.withTrackService(f, s.getMeTracks))
}

func (s *GinServer) getMeTracks(c *gin.Context, svc musichub.TrackService) {
	tracks, err := svc.FindMeTracks(c.Request.Context())
	if err != nil {
		s.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tracks)
}

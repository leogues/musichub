package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

func (s *GinServer) withArtistService(f provider.ArtistProviderFactoryInterface, handleFunc func(*gin.Context, musichub.ArtistService)) gin.HandlerFunc {
	return func(c *gin.Context) {
		providerName := c.Param("provider")

		token := musichub.ProviderTokenFromContext(c.Request.Context())

		service, err := f.CreateArtistProvider(providerName, token)

		if err != nil {
			s.Error(c, err)
			return
		}

		handleFunc(c, service)
	}
}

func (s *GinServer) registerArtistRoutes(r *gin.RouterGroup, f provider.ArtistProviderFactoryInterface) {
	r.GET("/me/platforms/:provider/artist", s.withArtistService(f, s.getMeArtists))
}

func (s *GinServer) getMeArtists(c *gin.Context, svc musichub.ArtistService) {

	artists, err := svc.FindMeArtists(c.Request.Context())

	if err != nil {
		s.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, artists)
}

package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

func (s *GinServer) withAlbumService(f provider.AlbumProviderFactoryInterface, handleFunc func(*gin.Context, musichub.AlbumService)) gin.HandlerFunc {
	return func(c *gin.Context) {
		providerName := c.Param("provider")

		token := musichub.ProviderTokenFromContext(c.Request.Context())

		service, err := f.CreateAlbumProvider(providerName, token)

		if err != nil {
			s.Error(c, err)
			return
		}

		handleFunc(c, service)
	}
}

func (s *GinServer) registerAlbumRoutes(r *gin.RouterGroup, f provider.AlbumProviderFactoryInterface) {
	r.GET("/me/platforms/:provider/album", s.withAlbumService(f, s.getMeAlbums))
	r.GET("/platforms/:provider/album/:id", s.withAlbumService(f, s.getAlbum))
}

func (s *GinServer) getMeAlbums(c *gin.Context, svc musichub.AlbumService) {

	albums, err := svc.FindMeAlbums(c.Request.Context())

	if err != nil {
		s.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, albums)
}

func (s *GinServer) getAlbum(c *gin.Context, svc musichub.AlbumService) {
	albumID := c.Param("id")

	album, err := svc.FindAlbumByID(c.Request.Context(), albumID)

	if err != nil {
		s.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, album)
}

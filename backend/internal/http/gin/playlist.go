package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

func (s *GinServer) withPlaylistService(f provider.PlaylistProviderFactoryInterface, handleFunc func(*gin.Context, musichub.PlaylistService)) gin.HandlerFunc {
	return func(c *gin.Context) {

		providerName := c.Param("provider")

		token := musichub.ProviderTokenFromContext(c.Request.Context())

		service, err := f.CreatePlaylistProvider(providerName, token)

		if err != nil {
			s.Error(c, err)
			return
		}

		handleFunc(c, service)
	}

}

func (s *GinServer) registerPlaylistRoutes(r *gin.RouterGroup, f provider.PlaylistProviderFactoryInterface) {
	r.GET("/me/platforms/:provider/playlist", s.withPlaylistService(f, s.getMePlaylists))
	r.GET("/platforms/:provider/playlist/:id", s.withPlaylistService(f, s.getPlaylist))
}

func (s *GinServer) getMePlaylists(c *gin.Context, svc musichub.PlaylistService) {
	playlists, err := svc.FindMePlaylists(c.Request.Context())

	if err != nil {
		s.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, playlists)
}

func (s *GinServer) getPlaylist(c *gin.Context, svc musichub.PlaylistService) {
	playlistID := c.Param("id")

	playlist, err := svc.FindPlaylistByID(c.Request.Context(), playlistID)

	if err != nil {
		s.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, playlist)
}

package gin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin/httpError"
	"github.com/leogues/MusicSyncHub/internal/http/gin/middleware"
	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
	"github.com/leogues/MusicSyncHub/internal/http/gin/strategy"
	"github.com/leogues/MusicSyncHub/internal/provider"
)

const ShutdownTimeout = 1 * time.Second

type GinServer struct {
	sessionManager  *session.SessionManager
	strategy        *strategy.OAuthStrategys
	strategyFactory *strategy.OAuthStrategyFactory
	middlware       *middleware.Middleware
	server          *http.Server
	Error           func(c *gin.Context, err error)

	IsProd   bool
	HashKey  string
	BlockKey string

	GoogleClientID     string
	GoogleClientSecret string

	SpotifyClientID     string
	SpotifyClientSecret string

	ProviderService     musichub.ProviderService
	UserService         musichub.UserService
	AuthService         musichub.AuthService
	ProviderAuthService musichub.ProviderAuthService

	ArtistProviderFactory   provider.ArtistProviderFactoryInterface
	AlbumProviderFactory    provider.AlbumProviderFactoryInterface
	TrackProviderFactory    provider.TrackProviderFactoryInterface
	PlaylistProviderFactory provider.PlaylistProviderFactoryInterface
}

func (s *GinServer) configureMiddleware() {
	s.middlware = middleware.NewMiddleware(s.sessionManager, s.strategyFactory, s.UserService, s.ProviderAuthService)
	s.middlware.Error = s.Error
}

func (s *GinServer) configureStrategys() {
	s.strategy = strategy.NewOAuthStrategys()
	s.strategy.Google = strategy.OAuth2GoogleConfig(s.GoogleClientID, s.GoogleClientSecret)
	s.strategy.Spotify = strategy.OAuth2SpotifyConfig(s.SpotifyClientID, s.SpotifyClientSecret)
	s.strategyFactory = strategy.NewOAuthStrategyFactory(s.strategy)
}

func (s *GinServer) Start() (*http.Server, error) {
	if s.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	sessionManager, err := session.NewSessionManager(s.HashKey, s.BlockKey)
	errorHandler := httpError.NewErrorHandler()

	if err != nil {
		return nil, err
	}

	s.configureStrategys()

	s.sessionManager = sessionManager
	s.Error = errorHandler

	s.configureMiddleware()

	r := gin.Default()
	r.Use(s.middlware.ReportPanic())
	r = s.configureHandlers(r)

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return s.server, nil
}

func (s *GinServer) configureHandlers(r *gin.Engine) *gin.Engine {

	r.Use(s.middlware.Authenticate())

	withNotRequiredAuth := r.Group("/api")
	withNotRequiredAuth.Use(s.middlware.RequireNoAuth())
	{
		s.registerAuthRoutes(withNotRequiredAuth)
	}

	withRequiredAuth := r.Group("/api")
	withRequiredAuth.Use(s.middlware.RequireAuth())
	{
		s.registerProviderAuthRoutes(withRequiredAuth)

		s.registerPlatformRoutes(withRequiredAuth)
		s.registerUserRoutes(withRequiredAuth)

		withMusicProvider := withRequiredAuth.Group("/")
		withMusicProvider.Use(s.middlware.RequireAuthProvider())
		{
			s.registerArtistRoutes(withMusicProvider, s.ArtistProviderFactory)
			s.registerAlbumRoutes(withMusicProvider, s.AlbumProviderFactory)
			s.registerTrackRoutes(withMusicProvider, s.TrackProviderFactory)
			s.registerPlaylistRoutes(withMusicProvider, s.PlaylistProviderFactory)
		}
	}

	return r
}

func (s *GinServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

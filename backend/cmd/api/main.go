package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leogues/MusicSyncHub/internal/http/gin"
	musichubapi "github.com/leogues/MusicSyncHub/internal/musichubAPI"
	"github.com/leogues/MusicSyncHub/internal/postgres"
	"github.com/leogues/MusicSyncHub/internal/provider"
	"github.com/leogues/MusicSyncHub/internal/sql"
)

type Main struct {
	Config     Config
	HTTPServer *gin.GinServer
	DB         *postgres.DB
}

func NewMain() *Main {
	return &Main{
		Config:     loadConfig(),
		HTTPServer: &gin.GinServer{},
		DB:         postgres.NewDB(""),
	}
}

func main() {
	// if os.Getenv("LOAD_ENV_FILE") == "true" {
	// 	if err := godotenv.Load(); err != nil {
	// 		panic(fmt.Errorf("ENV_FILE: %w", err))
	// 	}
	// }

	main := NewMain()

	if err := main.Run(); err != nil {
		log.Fatal("error running api", err)
		os.Exit(1)
	}
}

func (m *Main) Run() error {
	m.DB.DatasourceName = m.Config.DB.DatasourceName

	if err := m.DB.Open(); err != nil {
		panic(fmt.Errorf("database: %w", err))
	}

	sqlDB := m.DB.SqlDB()
	transactionManager := sql.NewTransaction(sqlDB)

	userRepository := postgres.NewUserRepository(m.DB)
	authRepository := postgres.NewAuthRepository(m.DB)
	providerAuthRepository := postgres.NewProviderAuthRepository(m.DB)

	userService := musichubapi.NewUserService(userRepository, authRepository, providerAuthRepository, transactionManager)
	authService := musichubapi.NewAuthService(authRepository, userRepository, transactionManager)
	providerAuthService := musichubapi.NewProviderAuth(providerAuthRepository, transactionManager)
	providerService := provider.NewProviderService()

	artistProviderFactory := provider.NewArtistProviderFactory()
	albumProviderFactory := provider.NewAlbumProviderFactory()
	trackProviderFactory := provider.NewTrackProviderFactory()
	playlistProviderFactory := provider.NewPlaylistProviderFactory()

	m.HTTPServer.IsProd = m.Config.HTTP.isProd
	m.HTTPServer.HashKey = m.Config.HTTP.HashKey
	m.HTTPServer.BlockKey = m.Config.HTTP.BlockKey

	m.HTTPServer.GoogleClientID = m.Config.Google.ClientID
	m.HTTPServer.GoogleClientSecret = m.Config.Google.ClientSecret
	m.HTTPServer.SpotifyClientID = m.Config.Spotify.ClientID
	m.HTTPServer.SpotifyClientSecret = m.Config.Spotify.ClientSecret

	m.HTTPServer.UserService = userService
	m.HTTPServer.AuthService = authService
	m.HTTPServer.ProviderAuthService = providerAuthService
	m.HTTPServer.ProviderService = providerService

	m.HTTPServer.ArtistProviderFactory = artistProviderFactory
	m.HTTPServer.AlbumProviderFactory = albumProviderFactory
	m.HTTPServer.TrackProviderFactory = trackProviderFactory
	m.HTTPServer.PlaylistProviderFactory = playlistProviderFactory

	server, err := m.HTTPServer.Start()

	if err != nil {
		log.Fatal("error starting server", err)
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

type HTTP struct {
	HashKey  string
	BlockKey string
	isProd   bool
}

type DB struct {
	DatasourceName string
}

type oAuth2Credentials struct {
	ClientID     string
	ClientSecret string
}

type Google oAuth2Credentials
type Spotify oAuth2Credentials

type Config struct {
	HTTP    HTTP
	DB      DB
	Google  Google
	Spotify Spotify
}

func newConfig() Config {
	return Config{}
}

func loadConfig() Config {
	config := newConfig()
	env := os.Getenv("ENV")
	config.HTTP.isProd = true

	if env == "development" {
		config.HTTP.isProd = false
	}

	config.DB.DatasourceName = os.Getenv("DATABASE_URL")

	config.HTTP.HashKey = os.Getenv("HASH_KEY")
	config.HTTP.BlockKey = os.Getenv("BLOCK_KEY")

	config.Google.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	config.Google.ClientSecret = os.Getenv("GOOGLE_SECRET")

	config.Spotify.ClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	config.Spotify.ClientSecret = os.Getenv("SPOTIFY_SECRET")

	return config
}

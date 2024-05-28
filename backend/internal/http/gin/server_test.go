package gin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin"
	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
	servicemock "github.com/leogues/MusicSyncHub/internal/mock/serviceMock"
)

const (
	TestHashKey             = "0000000000000000"
	TestHashBlock           = "00000000000000000000000000000000"
	TestGoogleClientID      = "GOOGLE_CLIENT_ID"
	TestGoogleClientSecret  = "GOOGLE_CLIENT_SECRET"
	TestSpotifyClientID     = "SPOTIFY_CLIENT_ID"
	TestSpotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
	BaseURL                 = "http://localhost:8080/api"
)

type Server struct {
	*gin.GinServer

	// Mocks
	ProviderService     servicemock.ProviderService
	UserService         servicemock.UserService
	AuthService         servicemock.AuthService
	ProviderAuthService servicemock.ProviderAuthService

	ArtistProviderFactory   servicemock.ArtistProviderFactory
	PlaylistProviderFactory servicemock.PlaylistProviderFactory
	TrackProviderFactory    servicemock.TrackProviderFactory
	AlbumProviderFactory    servicemock.AlbumProviderFactory

	sessionManager *session.SessionManager
}

func MustOpenServer(tb testing.TB) *Server {
	tb.Helper()
	server := &Server{GinServer: &gin.GinServer{}}

	server.HashKey = TestHashKey
	server.BlockKey = TestHashBlock
	server.GoogleClientID = TestGoogleClientID
	server.GoogleClientSecret = TestGoogleClientSecret
	server.SpotifyClientID = TestSpotifyClientID
	server.SpotifyClientSecret = TestSpotifyClientSecret

	server.GinServer.ProviderService = &server.ProviderService
	server.GinServer.UserService = &server.UserService
	server.GinServer.AuthService = &server.AuthService
	server.GinServer.ProviderAuthService = &server.ProviderAuthService

	server.GinServer.ArtistProviderFactory = &server.ArtistProviderFactory
	server.GinServer.PlaylistProviderFactory = &server.PlaylistProviderFactory
	server.GinServer.TrackProviderFactory = &server.TrackProviderFactory
	server.GinServer.AlbumProviderFactory = &server.AlbumProviderFactory
	var err error
	server.sessionManager, err = session.NewSessionManager(TestHashKey, TestHashBlock)
	if err != nil {
		tb.Fatalf("failed to create session manager: %v", err)
	}

	s, err := server.Start()
	if err != nil {
		tb.Fatalf("failed to start server: %v", err)
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				tb.Fatal(err)
			}
		}
	}()

	return server
}

func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}

func (s *Server) MustNewRequest(tb testing.TB, ctx context.Context, method, url string, body io.Reader) *http.Request {
	tb.Helper()

	r, err := http.NewRequest(method, BaseURL+url, nil)
	if err != nil {
		tb.Fatal(err)
	}

	if user := musichub.UserFromContext(ctx); user != nil {
		data, err := s.sessionManager.MarshalSession(session.Session{UserID: user.ID})
		if err != nil {
			tb.Fatal(err)
		}
		r.AddCookie(&http.Cookie{
			Name:  session.SessionCookieName,
			Value: data,
			Path:  "/",
		})
	}

	return r
}

func parseResponse[T any](tb testing.TB, response *http.Response) (T, error) {
	tb.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return *new(T), err
	}

	var result T

	err = json.Unmarshal(body, &result)
	if err != nil {
		return *new(T), err
	}

	return result, nil
}

func TestServer_Start(t *testing.T) {
	server := MustOpenServer(t)
	defer MustCloseServer(t, server)
}

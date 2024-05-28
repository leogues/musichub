package gin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	musichub "github.com/leogues/MusicSyncHub"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (s *GinServer) registerProviderAuthRoutes(r *gin.RouterGroup) {
	r.GET("/auth/spotify", s.oAuthSpotify)
	r.GET("/auth/spotify/callback", s.oAuthSpotifyCallback)
}

func (s *GinServer) oAuthSpotify(c *gin.Context) {
	session, err := s.sessionManager.Session(c)
	if err != nil {
		s.Error(c, err)
		return
	}

	state := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, state); err != nil {
		fmt.Print()
		s.Error(c, err)

		return
	}

	session.State = hex.EncodeToString(state)
	fmt.Print(session.State)

	if err := s.sessionManager.SetSession(c, session); err != nil {
		s.Error(c, err)
		return
	}

	c.Redirect(http.StatusFound, s.strategy.Spotify.AuthCodeURL(session.State, oauth2.AccessTypeOffline))
}

func (s *GinServer) oAuthSpotifyCallback(c *gin.Context) {
	state, code := c.Query("state"), c.Query("code")

	session, err := s.sessionManager.Session(c)
	if err != nil {
		s.Error(c, fmt.Errorf("cannot read session: %s", err))
		return
	}

	if state != session.State {
		s.Error(c, fmt.Errorf("oauth state mismatch"))
		return
	}

	user := musichub.UserFromContext(c.Request.Context())
	if user == nil {
		s.Error(c, musichub.Errorf(musichub.EUNAUTHORIZED, "Unauthorized access, User ID not found in context (%s)", musichub.AuthSourceSpotify))
		return
	}

	ctx := c.Request.Context()

	token, err := s.strategy.Spotify.Exchange(ctx, code)

	if err != nil {
		s.Error(c, err)
		return
	}

	providerAuth := &musichub.ProviderAuth{
		UserID:       user.ID,
		Source:       musichub.AuthSourceSpotify,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: &musichub.User{
			ID: user.ID,
		},
	}
	if !token.Expiry.IsZero() {
		providerAuth.Expiry = &token.Expiry
	}

	if err := s.ProviderAuthService.CreateProviderAuth(ctx, providerAuth); err != nil {
		s.Error(c, err)
		return
	}

	c.Redirect(http.StatusFound, "/authClosedAPI")
}

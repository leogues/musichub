package gin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin/strategy"
	"golang.org/x/oauth2"
)

func (s *GinServer) registerAuthRoutes(r *gin.RouterGroup) {
	r.GET("/auth/google", s.oAuthGoogle)
	r.GET("/auth/google/callback", s.oAuthGoogleCallback)

}

func (s *GinServer) oAuthGoogle(c *gin.Context) {
	redirectUrl := c.Query("redirect")

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

	session.RedirectURL = redirectUrl

	if err := s.sessionManager.SetSession(c, session); err != nil {
		s.Error(c, err)
		return
	}

	c.Redirect(http.StatusFound, s.strategy.Google.AuthCodeURL(session.State))
}

func (s *GinServer) oAuthGoogleCallback(c *gin.Context) {

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

	ctx := c.Request.Context()

	token, err := s.strategy.Google.Exchange(ctx, code, oauth2.AccessTypeOffline)

	if err != nil {
		s.Error(c, fmt.Errorf("oauth exchange error: %s", err.Error()))
		return

	}

	client := strategy.NewClientGoogle(token)

	user, err := client.GoogleUserInfo()

	if err != nil {
		s.Error(c, fmt.Errorf("cannot fetch user info: %s", err))

		return
	}

	if *user.ID == "" {
		s.Error(c, fmt.Errorf("user ID is not returned by Google, cannot authenticate user"))
		return
	}

	var name string
	if user.Name != nil {
		name = *user.Name
	}

	if user.FamilyName != nil {
		name += " " + *user.FamilyName
	}

	var email string

	if user.Email != nil {
		email = *user.Email
	}

	auth := &musichub.Auth{
		Source:       musichub.AuthSourceGoogle,
		SourceID:     *user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: &musichub.User{
			Name:  name,
			Email: email,
		},
	}
	if !token.Expiry.IsZero() {
		auth.Expiry = &token.Expiry
	}

	if err = s.AuthService.CreateAuth(ctx, auth); err != nil {
		s.Error(c, fmt.Errorf("cannot create auth: %s", err))
		return
	}

	redirectURL := session.RedirectURL

	session.UserID = auth.UserID
	session.RedirectURL = ""
	session.State = ""

	if err := s.sessionManager.SetSession(c, session); err != nil {
		s.Error(c, fmt.Errorf("cannot set session cookie: %s", err))
		return
	}

	if redirectURL == "" {
		redirectURL = "/"
	}

	c.Redirect(http.StatusFound, redirectURL)

}

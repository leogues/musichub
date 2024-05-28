package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"golang.org/x/oauth2"
)

func (m *Middleware) tokenExchange(ctx context.Context, providerAuth *musichub.ProviderAuth, providerName string) error {

	if providerAuth.RefreshToken == "" {
		return &musichub.Error{Code: musichub.EUNAUTHORIZED, Message: fmt.Sprintf("Unauthorized access to provider %s", providerName)}
	}

	strategy := m.strategyFactory.Strategy(ctx, providerName)
	token := &oauth2.Token{
		RefreshToken: providerAuth.RefreshToken,
	}

	tokenSource := strategy.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return err
	}

	providerAuth.AccessToken = newToken.AccessToken
	providerAuth.RefreshToken = newToken.RefreshToken
	providerAuth.Expiry = &newToken.Expiry

	go func() {
		if err := m.ProviderAuthService.UpdateProviderAuth(context.Background(), providerAuth); err != nil {
			log.Printf("cannot update provider auth: %s", err)
		}
	}()

	return nil
}

func (m *Middleware) RequireAuthProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		providerName := c.Param("provider")
		user := musichub.UserFromContext(ctx)

		var providerAuth *musichub.ProviderAuth

		for _, auth := range user.ProviderAuths {
			if auth.Source == providerName {
				providerAuth = auth
				break
			}
		}

		if providerAuth == nil {
			m.Error(c, musichub.Errorf(musichub.EUNAUTHORIZED, "Unauthorized access to provider %s", providerName))
			c.Abort()
			return
		}

		if providerAuth.Expiry.After(time.Now()) {
			c.Request = c.Request.WithContext(musichub.NewContextWithProviderToken(c.Request.Context(), providerAuth.AccessToken))
			c.Next()
			return
		}

		if err := m.tokenExchange(ctx, providerAuth, providerName); err != nil {
			m.Error(c, err)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(musichub.NewContextWithProviderToken(c.Request.Context(), providerAuth.AccessToken))

		c.Next()
	}
}

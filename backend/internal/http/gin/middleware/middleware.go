package middleware

import (
	"github.com/gin-gonic/gin"
	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/http/gin/session"
	"github.com/leogues/MusicSyncHub/internal/http/gin/strategy"
)

type Middleware struct {
	sessionManager  *session.SessionManager
	strategyFactory *strategy.OAuthStrategyFactory
	Error           func(c *gin.Context, err error)

	UserService         musichub.UserService
	ProviderAuthService musichub.ProviderAuthService
}

func NewMiddleware(sessionManager *session.SessionManager, strategyFactory *strategy.OAuthStrategyFactory, userService musichub.UserService, providerAuthService musichub.ProviderAuthService) *Middleware {
	return &Middleware{
		sessionManager:      sessionManager,
		strategyFactory:     strategyFactory,
		UserService:         userService,
		ProviderAuthService: providerAuthService,
	}
}

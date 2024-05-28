package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.AuthService = (*AuthService)(nil)

type AuthService struct {
	CreateAuthFn func(ctx context.Context, auth *musichub.Auth) error
}

func (s *AuthService) CreateAuth(ctx context.Context, auth *musichub.Auth) error {
	return s.CreateAuthFn(ctx, auth)
}

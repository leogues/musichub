package repositorymock

import (
	"context"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.AuthRepository = (*AuthRepository)(nil)

type AuthRepository struct {
	CreateAuthFn         func(ctx context.Context, auth *musichub.Auth) error
	UpdateAuthFn         func(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.Auth, error)
	FindAuthBySourceIDFn func(ctx context.Context, source, sourceID string) (*musichub.Auth, error)
	FindAuthByIDFn       func(ctx context.Context, id int) (*musichub.Auth, error)
	AttachUserAuthsFn    func(ctx context.Context, user *musichub.User) error
}

func (r *AuthRepository) CreateAuth(ctx context.Context, auth *musichub.Auth) error {
	return r.CreateAuthFn(ctx, auth)
}

func (r *AuthRepository) UpdateAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.Auth, error) {
	return r.UpdateAuthFn(ctx, id, accessToken, refreshToken, expiry)
}

func (r *AuthRepository) FindAuthBySourceID(ctx context.Context, source, sourceID string) (*musichub.Auth, error) {
	return r.FindAuthBySourceIDFn(ctx, source, sourceID)
}

func (r *AuthRepository) FindAuthByID(ctx context.Context, id int) (*musichub.Auth, error) {
	return r.FindAuthByIDFn(ctx, id)
}

func (r *AuthRepository) AttachUserAuths(ctx context.Context, user *musichub.User) error {
	return r.AttachUserAuthsFn(ctx, user)
}

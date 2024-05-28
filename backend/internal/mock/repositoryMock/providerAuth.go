package repositorymock

import (
	"context"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.ProviderAuthRepository = (*ProviderAuthRepository)(nil)

type ProviderAuthRepository struct {
	CreateProviderAuthFn       func(ctx context.Context, providerAuth *musichub.ProviderAuth) error
	UpdateProviderAuthFn       func(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.ProviderAuth, error)
	FindProviderAuthBySourceFn func(ctx context.Context, id int) (*musichub.ProviderAuth, error)
	AttachUserProviderAuthsFn  func(ctx context.Context, user *musichub.User) error
}

func (r *ProviderAuthRepository) CreateProviderAuth(ctx context.Context, auth *musichub.ProviderAuth) error {
	return r.CreateProviderAuthFn(ctx, auth)
}

func (r *ProviderAuthRepository) UpdateProviderAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.ProviderAuth, error) {
	return r.UpdateProviderAuthFn(ctx, id, accessToken, refreshToken, expiry)
}

func (r *ProviderAuthRepository) FindProviderAuthBySource(ctx context.Context, userId int, source string) (*musichub.ProviderAuth, error) {
	return r.FindProviderAuthBySourceFn(ctx, userId)
}

func (r *ProviderAuthRepository) AttachUserProviderAuths(ctx context.Context, user *musichub.User) error {
	return r.AttachUserProviderAuthsFn(ctx, user)
}

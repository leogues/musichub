package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.ProviderAuthService = (*ProviderAuthService)(nil)

type ProviderAuthService struct {
	CreateProviderAuthFn       func(context.Context, *musichub.ProviderAuth) error
	UpdateProviderAuthFn       func(context.Context, *musichub.ProviderAuth) error
	FindProviderAuthBySourceFn func(context.Context, int, string) (*musichub.ProviderAuth, error)
}

func (s *ProviderAuthService) CreateProviderAuth(ctx context.Context, auth *musichub.ProviderAuth) error {
	return s.CreateProviderAuthFn(ctx, auth)
}

func (s *ProviderAuthService) UpdateProviderAuth(ctx context.Context, auth *musichub.ProviderAuth) error {
	return s.UpdateProviderAuthFn(ctx, auth)
}

func (s *ProviderAuthService) FindProviderAuthBySource(ctx context.Context, source int, sourceID string) (*musichub.ProviderAuth, error) {
	return s.FindProviderAuthBySourceFn(ctx, source, sourceID)
}

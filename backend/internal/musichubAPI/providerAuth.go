package musichubapi

import (
	"context"
	"fmt"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.ProviderAuthService = (*ProviderAuthService)(nil)

type ProviderAuthService struct {
	providerAuthRepository musichub.ProviderAuthRepository
	tx                     musichub.Transaction
}

func NewProviderAuth(providerAuthRepository musichub.ProviderAuthRepository, tx musichub.Transaction) *ProviderAuthService {
	return &ProviderAuthService{
		providerAuthRepository: providerAuthRepository,
		tx:                     tx,
	}
}

func (s *ProviderAuthService) FindProviderAuthBySource(ctx context.Context, userID int, source string) (*musichub.ProviderAuth, error) {
	providerAuth, err := s.providerAuthRepository.FindProviderAuthBySource(ctx, userID, source)

	if err != nil && musichub.ErrorCode(err) != musichub.ENOTFOUND {
		return nil, fmt.Errorf("cannot find provider_auth by source user: %w", err)
	}

	if err != nil {
		return nil, err
	}

	return providerAuth, nil
}

func (s *ProviderAuthService) UpdateProviderAuth(ctx context.Context, providerAuth *musichub.ProviderAuth) error {
	updatedProviderAuth, err := s.providerAuthRepository.UpdateProviderAuth(ctx, providerAuth.ID, providerAuth.AccessToken, providerAuth.RefreshToken, providerAuth.Expiry)

	if err != nil {
		return fmt.Errorf("connot update provider_auth: id=%d err=%w", providerAuth.ID, err)
	}

	*providerAuth = *updatedProviderAuth

	return nil
}

func (s *ProviderAuthService) CreateProviderAuth(ctx context.Context, providerAuth *musichub.ProviderAuth) error {
	if err := s.tx.Begin(&ctx); err != nil {
		return err
	}
	defer s.tx.Rollback(ctx)

	existingProviderAuth, err := s.providerAuthRepository.FindProviderAuthBySource(ctx, providerAuth.UserID, providerAuth.Source)

	if err != nil && musichub.ErrorCode(err) != musichub.ENOTFOUND {
		return fmt.Errorf("cannot find provider_auth by source user: %w", err)
	}

	if existingProviderAuth != nil {
		if existingProviderAuth, err = s.providerAuthRepository.UpdateProviderAuth(ctx, existingProviderAuth.ID, providerAuth.AccessToken, providerAuth.RefreshToken, providerAuth.Expiry); err != nil {
			return fmt.Errorf("connot update provider_auth: id=%d err=%w", existingProviderAuth.ID, err)
		}
		*providerAuth = *existingProviderAuth
		return s.tx.Commit(ctx)
	}

	if err = s.providerAuthRepository.CreateProviderAuth(ctx, providerAuth); err != nil {
		return err
	}

	return s.tx.Commit(ctx)
}

package musichubapi

import (
	"context"
	"fmt"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.AuthService = (*AuthService)(nil)

type AuthService struct {
	authRepository musichub.AuthRepository
	userRepository musichub.UserRepository
	tx             musichub.Transaction
}

func NewAuthService(authRepository musichub.AuthRepository, userRepository musichub.UserRepository, tx musichub.Transaction) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		userRepository: userRepository,
		tx:             tx,
	}
}

func (s *AuthService) CreateAuth(ctx context.Context, auth *musichub.Auth) error {
	if err := s.tx.Begin(&ctx); err != nil {
		return err
	}
	defer s.tx.Rollback(ctx)

	existingAuth, err := s.authRepository.FindAuthBySourceID(ctx, auth.Source, auth.SourceID)

	if err != nil && musichub.ErrorCode(err) != musichub.ENOTFOUND {
		return fmt.Errorf("cannot find auth by source user: %w", err)
	}

	if existingAuth != nil {
		if existingAuth, err = s.authRepository.UpdateAuth(ctx, existingAuth.ID, auth.AccessToken, auth.RefreshToken, auth.Expiry); err != nil {
			return fmt.Errorf("connot update auth: id=%d err=%w", existingAuth.ID, err)
		}
		*auth = *existingAuth
		return s.tx.Commit(ctx)
	}

	if auth.User != nil && auth.UserID == 0 {
		if err = s.createUserIfNeeded(ctx, auth); err != nil {
			return err
		}
	}

	if err = s.authRepository.CreateAuth(ctx, auth); err != nil {
		return err
	}

	if err = s.userRepository.AttachAuthAssociations(ctx, auth); err != nil {
		return err
	}

	return s.tx.Commit(ctx)

}

func (s *AuthService) createUserIfNeeded(ctx context.Context, auth *musichub.Auth) error {
	user, err := s.userRepository.FindUserByEmail(ctx, auth.User.Email)

	if err != nil && musichub.ErrorCode(err) != musichub.ENOTFOUND {
		return fmt.Errorf("cannot find user by email: %w", err)
	}

	if user != nil {
		auth.User = user
		auth.UserID = auth.User.ID
		return nil
	}

	if err = s.userRepository.CreateUser(ctx, auth.User); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	auth.UserID = auth.User.ID
	return nil
}

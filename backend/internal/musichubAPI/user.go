package musichubapi

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.UserService = (*UserService)(nil)

type UserService struct {
	userRepository         musichub.UserRepository
	authRepository         musichub.AuthRepository
	providerAuthRepository musichub.ProviderAuthRepository
	tx                     musichub.Transaction
}

func NewUserService(userRepository musichub.UserRepository, authRepository musichub.AuthRepository, providerAuthRepository musichub.ProviderAuthRepository, tx musichub.Transaction) *UserService {
	return &UserService{
		userRepository:         userRepository,
		authRepository:         authRepository,
		providerAuthRepository: providerAuthRepository,
		tx:                     tx,
	}
}

func (s *UserService) FindUsers(ctx context.Context, filter musichub.UserFilter) ([]*musichub.User, int, error) {
	if err := s.tx.Begin(&ctx); err != nil {
		return nil, 0, err
	}
	defer s.tx.Rollback(ctx)

	users, n, err := s.userRepository.FindUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if err = s.tx.Commit(ctx); err != nil {
		return nil, 0, err
	}

	return users, n, nil

}

func (s *UserService) FindUserByID(ctx context.Context, id int) (*musichub.User, error) {

	if err := s.tx.Begin(&ctx); err != nil {
		return nil, err
	}
	defer s.tx.Rollback(ctx)

	user, err := s.userRepository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	} else if err := s.authRepository.AttachUserAuths(ctx, user); err != nil {
		return user, err
	} else if err := s.providerAuthRepository.AttachUserProviderAuths(ctx, user); err != nil {
		return user, err
	}

	if err = s.tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) CreateUser(ctx context.Context, user *musichub.User) error {
	if err := s.tx.Begin(&ctx); err != nil {
		return err
	}
	defer s.tx.Rollback(ctx)

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return err
	}

	return s.tx.Commit(ctx)

}

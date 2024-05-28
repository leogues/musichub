package servicemock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.UserService = (*UserService)(nil)

type UserService struct {
	FindUsersByIDFn func(context.Context, int) (*musichub.User, error)
	FindUsersFn     func(context.Context, musichub.UserFilter) ([]*musichub.User, int, error)
	CreateUserFn    func(context.Context, *musichub.User) error
}

func (s *UserService) FindUserByID(ctx context.Context, id int) (*musichub.User, error) {
	return s.FindUsersByIDFn(ctx, id)
}

func (s *UserService) FindUsers(ctx context.Context, filter musichub.UserFilter) ([]*musichub.User, int, error) {
	return s.FindUsersFn(ctx, filter)
}

func (s *UserService) CreateUser(ctx context.Context, user *musichub.User) error {
	return s.CreateUserFn(ctx, user)
}

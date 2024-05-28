package repositorymock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	CreateUserFn             func(ctx context.Context, user *musichub.User) error
	FindUserByIDFn           func(ctx context.Context, id int) (*musichub.User, error)
	FindUserByEmailFn        func(ctx context.Context, email string) (*musichub.User, error)
	FindUsersFn              func(ctx context.Context, filter musichub.UserFilter) ([]*musichub.User, int, error)
	AttachAuthAssociationsFn func(ctx context.Context, auth *musichub.Auth) error
}

func (r *UserRepository) CreateUser(ctx context.Context, user *musichub.User) error {
	return r.CreateUserFn(ctx, user)
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*musichub.User, error) {
	return r.FindUserByIDFn(ctx, id)
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*musichub.User, error) {
	return r.FindUserByEmailFn(ctx, email)
}

func (r *UserRepository) FindUsers(ctx context.Context, filter musichub.UserFilter) ([]*musichub.User, int, error) {
	return r.FindUsersFn(ctx, filter)
}

func (r *UserRepository) AttachAuthAssociations(ctx context.Context, auth *musichub.Auth) error {
	return r.AttachAuthAssociationsFn(ctx, auth)
}

package musichub

import (
	"context"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	APIKey string `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Auths         []*Auth         `json:"auths"`
	ProviderAuths []*ProviderAuth `json:"provider_auths"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return Errorf(EINVALID, "User name is required")
	}

	return nil
}

type UserWriter interface {
	CreateUser(ctx context.Context, user *User) error
}

type UserReader interface {
	FindUserByID(ctx context.Context, id int) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)
	AttachAuthAssociations(ctx context.Context, auth *Auth) error
}

type UserRepository interface {
	UserWriter
	UserReader
}

type UserService interface {
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)
	FindUserByID(ctx context.Context, id int) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

type UserFilter struct {
	ID     *int    `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

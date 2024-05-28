package musichub

import (
	"context"
	"time"
)

type ProviderAuth struct {
	ID int `json:"id"`

	UserID int   `json:"-"`
	User   *User `json:"-"`

	Source string `json:"source"`

	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	Expiry       *time.Time `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (a *ProviderAuth) Validate() error {
	if a.UserID == 0 {
		return Errorf(EINVALID, "User required.")
	} else if a.Source == "" {
		return Errorf(EINVALID, "Source required.")
	} else if a.AccessToken == "" {
		return Errorf(EINVALID, "Access token required.")
	}
	return nil
}

type ProviderAuthWriter interface {
	CreateProviderAuth(ctx context.Context, auth *ProviderAuth) error
	UpdateProviderAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*ProviderAuth, error)
}

type ProviderAuthReader interface {
	FindProviderAuthBySource(ctx context.Context, userId int, source string) (*ProviderAuth, error)
	AttachUserProviderAuths(ctx context.Context, user *User) error
}

type ProviderAuthRepository interface {
	ProviderAuthWriter
	ProviderAuthReader
}

type ProviderAuthService interface {
	CreateProviderAuth(ctx context.Context, auth *ProviderAuth) error
	UpdateProviderAuth(ctx context.Context, auth *ProviderAuth) error
	FindProviderAuthBySource(ctx context.Context, userID int, source string) (*ProviderAuth, error)
}

type ProviderAuthFilter struct {
	ID     *int    `json:"id"`
	UserID *int    `json:"userID"`
	Source *string `json:"source"`
}

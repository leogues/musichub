package musichub

import (
	"context"
	"time"
)

const (
	AuthSourceGoogle  = "google"
	AuthSourceSpotify = "spotify"
)

type Auth struct {
	ID int `json:"id"`

	UserID int   `json:"-"`
	User   *User `json:"-"`

	Source   string `json:"source"`
	SourceID string `json:"source_id"`

	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	Expiry       *time.Time `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (a *Auth) Validate() error {
	if a.UserID == 0 {
		return Errorf(EINVALID, "User required.")
	} else if a.Source == "" {
		return Errorf(EINVALID, "Source required.")
	} else if a.SourceID == "" {
		return Errorf(EINVALID, "Source ID required.")
	} else if a.AccessToken == "" {
		return Errorf(EINVALID, "Access token required.")
	}
	return nil
}

type AuthWriter interface {
	CreateAuth(ctx context.Context, auth *Auth) error
	UpdateAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*Auth, error)
}

type AuthReader interface {
	FindAuthBySourceID(ctx context.Context, source, sourceID string) (*Auth, error)
	FindAuthByID(ctx context.Context, id int) (*Auth, error)
	AttachUserAuths(ctx context.Context, user *User) error
}

type AuthRepository interface {
	AuthWriter
	AuthReader
}

type AuthService interface {
	CreateAuth(ctx context.Context, auth *Auth) error
}

type AuthFilter struct {
	ID       *int    `json:"id"`
	UserID   *int    `json:"userID"`
	Source   *string `json:"source"`
	SourceID *string `json:"sourceID"`
}

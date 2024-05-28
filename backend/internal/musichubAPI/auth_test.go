package musichubapi

import (
	"context"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	repositorymock "github.com/leogues/MusicSyncHub/internal/mock/repositoryMock"
)

func TestAuthService_CreateAuth(t *testing.T) {
	authRepository := &repositorymock.AuthRepository{}
	userRepository := &repositorymock.UserRepository{}
	tx := &repositorymock.Transaction{}

	auth0 := &musichub.Auth{
		ID:          1,
		UserID:      1,
		Source:      musichub.AuthSourceGoogle,
		SourceID:    "SOURCEID",
		AccessToken: "ACCESSTOKEN",
	}

	user0 := &musichub.User{
		ID:        1,
		Name:      "jhon",
		Email:     "jhon@gmail.com",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("OKCreateAuth", func(t *testing.T) {
		authService := NewAuthService(authRepository, userRepository, tx)

		authRepository.FindAuthBySourceIDFn = func(ctx context.Context, source string, sourceID string) (*musichub.Auth, error) {
			return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "Auth not found."}
		}

		userRepository.FindUserByEmailFn = func(ctx context.Context, email string) (*musichub.User, error) {
			return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "User not found."}
		}

		userRepository.CreateUserFn = func(ctx context.Context, user *musichub.User) error {
			*user = *user0
			return nil
		}

		authRepository.CreateAuthFn = func(ctx context.Context, auth *musichub.Auth) error {
			*auth = *auth0
			return nil
		}

		userRepository.AttachAuthAssociationsFn = func(ctx context.Context, auth *musichub.Auth) error {
			auth.User = user0
			return nil
		}

		auth := &musichub.Auth{
			User: &musichub.User{},
		}

		if err := authService.CreateAuth(context.Background(), auth); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		user := auth.User
		if !reflect.DeepEqual(user0, user) {
			t.Fatalf("expected %v, got %v", user0, user)
		}

		auth.User = nil
		if !reflect.DeepEqual(auth, auth0) {
			t.Fatalf("expected %v, got %v", auth0, auth)
		}
	})

	t.Run("OKCreateWithUserExists", func(t *testing.T) {
		authService := NewAuthService(authRepository, userRepository, tx)

		authRepository.FindAuthBySourceIDFn = func(ctx context.Context, source string, sourceID string) (*musichub.Auth, error) {
			return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "Auth not found."}
		}

		userRepository.FindUserByEmailFn = func(ctx context.Context, email string) (*musichub.User, error) {
			return user0, nil
		}

		authRepository.CreateAuthFn = func(ctx context.Context, auth *musichub.Auth) error {
			*auth = *auth0
			return nil
		}

		userRepository.AttachAuthAssociationsFn = func(ctx context.Context, auth *musichub.Auth) error {
			auth.User = user0
			return nil
		}

		auth := &musichub.Auth{
			User: &musichub.User{},
		}

		err := authService.CreateAuth(context.Background(), auth)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		user := auth.User
		if !reflect.DeepEqual(user0, user) {
			t.Fatalf("expected %v, got %v", user0, user)
		}

		auth.User = nil
		if !reflect.DeepEqual(auth, auth0) {
			t.Fatalf("expected %v, got %v", auth0, auth)
		}
	})

	t.Run("OKAuthExists", func(t *testing.T) {
		authService := NewAuthService(authRepository, userRepository, tx)

		authRepository.FindAuthBySourceIDFn = func(ctx context.Context, source string, sourceID string) (*musichub.Auth, error) {
			return auth0, nil
		}

		authRepository.UpdateAuthFn = func(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.Auth, error) {
			auth0.AccessToken = accessToken
			return auth0, nil
		}

		auth := &musichub.Auth{
			User: &musichub.User{},
		}

		err := authService.CreateAuth(context.Background(), auth)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(auth, auth0) {
			t.Fatalf("expected %v, got %v", auth0, auth)
		}
	})

}

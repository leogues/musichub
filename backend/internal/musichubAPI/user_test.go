package musichubapi

import (
	"context"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	repositorymock "github.com/leogues/MusicSyncHub/internal/mock/repositoryMock"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepository := &repositorymock.UserRepository{}
	tx := &repositorymock.Transaction{}

	user0 := &musichub.User{
		ID:        1,
		Name:      "jhon",
		Email:     "jhon@gmail.com",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	t.Run("OKCreateUser", func(t *testing.T) {
		userService := NewUserService(userRepository, nil, nil, tx)

		userRepository.CreateUserFn = func(ctx context.Context, user *musichub.User) error {
			*user = *user0
			return nil
		}

		user := &musichub.User{}

		err := userService.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(user0, user) {
			t.Fatalf("mismatch: %#v != %#v", user0, user)
		}
	})
}

func TestUserService_FindUserByID(t *testing.T) {
	userRepository := &repositorymock.UserRepository{}
	authRepository := &repositorymock.AuthRepository{}
	providerAuthRepository := &repositorymock.ProviderAuthRepository{}
	tx := &repositorymock.Transaction{}

	auth0 := &musichub.Auth{
		ID:          1,
		UserID:      1,
		Source:      musichub.AuthSourceGoogle,
		SourceID:    "SOURCEID",
		AccessToken: "ACCESSTOKEN",
	}

	providerAuth0 := &musichub.ProviderAuth{
		ID:          1,
		UserID:      1,
		Source:      musichub.AuthSourceSpotify,
		AccessToken: "ACCESSTOKEN",
	}

	user0 := &musichub.User{
		ID:    1,
		Name:  "jhon",
		Email: "jhon@gmail.com",
		Auths: []*musichub.Auth{
			auth0,
		},
		ProviderAuths: []*musichub.ProviderAuth{
			providerAuth0,
		},
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("OKFindUserByID", func(t *testing.T) {
		userService := NewUserService(userRepository, authRepository, providerAuthRepository, tx)

		userRepository.FindUserByIDFn = func(ctx context.Context, id int) (*musichub.User, error) {
			user := user0
			user.Auths = []*musichub.Auth{}
			user.ProviderAuths = []*musichub.ProviderAuth{}

			return user, nil
		}

		authRepository.AttachUserAuthsFn = func(ctx context.Context, user *musichub.User) error {
			user.Auths = []*musichub.Auth{auth0}
			return nil
		}

		providerAuthRepository.AttachUserProviderAuthsFn = func(ctx context.Context, user *musichub.User) error {
			user.ProviderAuths = []*musichub.ProviderAuth{providerAuth0}
			return nil
		}

		other, err := userService.FindUserByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		} else if !reflect.DeepEqual(user0, other) {
			t.Fatalf("mismatch: %#v != %#v", user0, other)
		}
	})

}

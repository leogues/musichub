package musichubapi

import (
	"context"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	repositorymock "github.com/leogues/MusicSyncHub/internal/mock/repositoryMock"
)

func TestProviderAuthService_CreateProviderAuth(t *testing.T) {
	providerAuthRepository := &repositorymock.ProviderAuthRepository{}
	tx := &repositorymock.Transaction{}

	providerAuth0 := &musichub.ProviderAuth{
		ID:           1,
		UserID:       1,
		Source:       musichub.AuthSourceSpotify,
		AccessToken:  "ACCESSTOKEN",
		RefreshToken: "REFRESHTOKEN",
	}

	t.Run("OKCreateProviderAuth", func(t *testing.T) {
		providerAuthService := NewProviderAuth(providerAuthRepository, tx)

		providerAuthRepository.FindProviderAuthBySourceFn = func(ctx context.Context, id int) (*musichub.ProviderAuth, error) {
			return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "ProviderAuth not found."}
		}

		providerAuthRepository.CreateProviderAuthFn = func(ctx context.Context, providerAuth *musichub.ProviderAuth) error {
			*providerAuth = *providerAuth0
			return nil
		}

		providerAuth := &musichub.ProviderAuth{}
		if err := providerAuthService.CreateProviderAuth(context.Background(), providerAuth); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(providerAuth0, providerAuth) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth0, providerAuth)
		}
	})

	t.Run("OKProviderAuthExists", func(t *testing.T) {
		providerAuthService := NewProviderAuth(providerAuthRepository, tx)

		providerAuthRepository.FindProviderAuthBySourceFn = func(ctx context.Context, id int) (*musichub.ProviderAuth, error) {
			return providerAuth0, nil
		}

		updatedProviderAuth0 := *providerAuth0
		updatedProviderAuth0.AccessToken = "UPDATEDACCESSTOKEN"

		providerAuthRepository.UpdateProviderAuthFn = func(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.ProviderAuth, error) {
			return &updatedProviderAuth0, nil
		}

		providerAuth := &musichub.ProviderAuth{}

		err := providerAuthService.CreateProviderAuth(context.Background(), providerAuth)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(providerAuth, &updatedProviderAuth0) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth, updatedProviderAuth0)
		}
	})
}

func TestProviderAuthService_UpdateProviderAuth(t *testing.T) {
	providerAuthRepository := &repositorymock.ProviderAuthRepository{}
	tx := &repositorymock.Transaction{}

	providerAuth0 := &musichub.ProviderAuth{
		ID:           1,
		UserID:       1,
		Source:       musichub.AuthSourceSpotify,
		AccessToken:  "ACCESSTOKEN",
		RefreshToken: "REFRESHTOKEN",
	}

	t.Run("OKUpdateProviderAuth", func(t *testing.T) {
		providerAuthService := NewProviderAuth(providerAuthRepository, tx)

		updatedProviderAuth0 := *providerAuth0
		updatedProviderAuth0.AccessToken = "UPDATEDACCESSTOKEN"

		providerAuthRepository.UpdateProviderAuthFn = func(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.ProviderAuth, error) {
			return &updatedProviderAuth0, nil
		}

		providerAuth := &musichub.ProviderAuth{}

		err := providerAuthService.UpdateProviderAuth(context.Background(), providerAuth)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(providerAuth, &updatedProviderAuth0) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth, updatedProviderAuth0)
		}

	})
}

func TestProviderAuthService_FindProviderAuthBySource(t *testing.T) {
	providerAuthRepository := &repositorymock.ProviderAuthRepository{}
	tx := &repositorymock.Transaction{}

	t.Run("OKFindProviderAuthBySource", func(t *testing.T) {
		providerAuthService := NewProviderAuth(providerAuthRepository, tx)

		providerAuthRepository.FindProviderAuthBySourceFn = func(ctx context.Context, id int) (*musichub.ProviderAuth, error) {
			return &musichub.ProviderAuth{}, nil
		}

		providerAuth, err := providerAuthService.FindProviderAuthBySource(context.Background(), 1, "source")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if providerAuth == nil {
			t.Fatal("expected not nil, got nil")
		}
	})
}

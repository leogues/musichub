package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/postgres"
)

func TestProviderAuthRepository_CreateProviderAuth(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)

	user := &musichub.User{
		Name:  "jhon",
		Email: "jhon@gmail.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	t.Run("OKCreateProviderAuth", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		expiry := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

		providerAuth := &musichub.ProviderAuth{
			Source:       musichub.AuthSourceSpotify,
			AccessToken:  "ACCESS_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
			Expiry:       &expiry,
			User:         user,
			UserID:       user.ID,
		}

		if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err != nil {
			t.Fatal(err)
		} else if got, want := providerAuth.ID, 1; got != want {
			t.Fatalf("got %v, want %v", got, want)
		} else if providerAuth.CreatedAt.IsZero() {
			t.Fatal("expected created at")
		} else if providerAuth.UpdatedAt.IsZero() {
			t.Fatal("expected updated at")
		}

		providerAuth.User = nil

		if other, err := providerAuthRepository.FindProviderAuthByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(providerAuth, other) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth, other)
		}
	})

	t.Run("ErrUserRequired", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		providerAuth := &musichub.ProviderAuth{
			Source:       musichub.AuthSourceGoogle,
			AccessToken:  "ACCESS_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
		}

		if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "User required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("ErrSourceRequired", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		providerAuth := &musichub.ProviderAuth{
			User:        user,
			UserID:      user.ID,
			AccessToken: "ACCESS_TOKEN",
		}

		if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "Source required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("ErrAccessTokenRequired", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		providerAuth := &musichub.ProviderAuth{
			User:   user,
			UserID: user.ID,
			Source: musichub.AuthSourceGoogle,
		}

		if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "Access token required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestAuthRepository_UpdateProviderAuth(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)

	user := &musichub.User{
		Name:  "jhon",
		Email: "jhon@gmail.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	t.Run("OKUpdateProviderAuth", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		expiry := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		providerAuth := &musichub.ProviderAuth{
			Source:       musichub.AuthSourceSpotify,
			AccessToken:  "ACCESS_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
			Expiry:       &expiry,
			User:         user,
			UserID:       user.ID,
		}
		if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err != nil {
			t.Fatal(err)
		}

		providerAuth.User = nil

		newAcessToken, newRefreshToken, newExpiry := "NEW_ACCESS_TOKEN", "NEW_REFRESH_TOKEN", time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
		updatedProviderAuth, err := providerAuthRepository.UpdateProviderAuth(context.Background(), 1, newAcessToken, newRefreshToken, &newExpiry)

		if err != nil {
			t.Fatal(err)
		} else if got, want := updatedProviderAuth.AccessToken, newAcessToken; got != want {
			t.Fatalf("AccessToken=%v, want %v", got, want)
		} else if got, want := updatedProviderAuth.RefreshToken, newRefreshToken; got != want {
			t.Fatalf("RefreshToken=%v, want %v", got, want)
		} else if got, want := updatedProviderAuth.Expiry, &newExpiry; !reflect.DeepEqual(got, want) {
			t.Fatalf("Expiry=%v, want %v", got, want)
		}

		if other, err := providerAuthRepository.FindProviderAuthByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(updatedProviderAuth, other) {
			t.Fatalf("mismatch: %#v != %#v", updatedProviderAuth, other)
		}
	})
}

func TestProviderAuthRepository_FindAuth(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)
	providerAuthRepository := postgres.NewProviderAuthRepository(db)

	user := &musichub.User{
		Name:  "jhon",
		Email: "jhon@gmai.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	providerAuth := &musichub.ProviderAuth{
		Source:       musichub.AuthSourceSpotify,
		AccessToken:  "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
		User:         user,
		UserID:       user.ID,
	}

	if err := providerAuthRepository.CreateProviderAuth(context.Background(), providerAuth); err != nil {
		t.Fatal(err)
	}

	providerAuth.User = nil

	t.Run("OKFindProviderAuthByID", func(t *testing.T) {
		if other, err := providerAuthRepository.FindProviderAuthByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(providerAuth, other) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth, other)
		}
	})

	t.Run("OKFindProviderAuthBySource", func(t *testing.T) {
		if other, err := providerAuthRepository.FindProviderAuthBySource(context.Background(), 1, providerAuth.Source); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(providerAuth, other) {
			t.Fatalf("mismatch: %#v != %#v", providerAuth, other)
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		if _, err := providerAuthRepository.FindProviderAuthByID(context.Background(), 2); musichub.ErrorCode(err) != musichub.ENOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestProviderAuthRepository_AttachUserProviderAuths(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)

	user := &musichub.User{
		Name:  "jhon",
		Email: "jhon@gmail.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	t.Run("OKAttachUserProviderAuths", func(t *testing.T) {
		providerAuthRepository := postgres.NewProviderAuthRepository(db)

		if err := providerAuthRepository.AttachUserProviderAuths(context.Background(), user); err != nil {
			t.Fatal(err)
		}

		if user.ProviderAuths == nil {
			t.Fatal("expected provider auths")
		}
	})

}

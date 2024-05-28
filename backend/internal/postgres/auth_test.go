package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/postgres"
)

func TestAuthRepository_CreateAuth(t *testing.T) {
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

	t.Run("OKCreateAuth", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)

		expiry := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		auth := &musichub.Auth{
			Source:       musichub.AuthSourceGoogle,
			SourceID:     "SOURCEID",
			AccessToken:  "ACCESS_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
			Expiry:       &expiry,
			User:         user,
			UserID:       user.ID,
		}

		if err := authRepository.CreateAuth(context.Background(), auth); err != nil {
			t.Fatal(err)
		} else if got, want := auth.ID, 1; got != want {
			t.Fatalf("got %v, want %v", got, want)
		} else if auth.CreatedAt.IsZero() {
			t.Fatal("expected created at")
		} else if auth.UpdatedAt.IsZero() {
			t.Fatal("expected updated at")
		}

		auth.User = nil

		if other, err := authRepository.FindAuthByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(auth, other) {
			t.Fatalf("mismatch: %#v != %#v", auth, other)
		}
	})

	t.Run("ErrUserRequired", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)
		auth := &musichub.Auth{
			Source:      musichub.AuthSourceGoogle,
			SourceID:    "SOURCEID",
			AccessToken: "ACCESS_TOKEN",
		}
		if err := authRepository.CreateAuth(context.Background(), auth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "User required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("ErrSourceRequired", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)
		auth := &musichub.Auth{
			User:        user,
			UserID:      user.ID,
			SourceID:    "SOURCEID",
			AccessToken: "ACCESS_TOKEN",
		}
		if err := authRepository.CreateAuth(context.Background(), auth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "Source required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("ErrSourceIDRequired", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)
		auth := &musichub.Auth{
			User:        user,
			UserID:      user.ID,
			Source:      musichub.AuthSourceGoogle,
			AccessToken: "ACCESS_TOKEN",
		}
		if err := authRepository.CreateAuth(context.Background(), auth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "Source ID required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("ErrAccessTokenRequired", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)
		auth := &musichub.Auth{
			User:     user,
			UserID:   user.ID,
			Source:   musichub.AuthSourceGoogle,
			SourceID: "SOURCEID",
		}
		if err := authRepository.CreateAuth(context.Background(), auth); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "Access token required." {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestAuthRepository_UpdateAuth(t *testing.T) {
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
	t.Run("OKUpdateAuth", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)

		expiry := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		auth := &musichub.Auth{
			Source:       musichub.AuthSourceGoogle,
			SourceID:     "SOURCEID",
			AccessToken:  "ACCESS_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
			Expiry:       &expiry,
			User:         user,
			UserID:       user.ID,
		}
		if err := authRepository.CreateAuth(context.Background(), auth); err != nil {
			t.Fatal(err)
		}

		auth.User = nil

		newAcessToken, newRefreshToken, newExpiry := "NEW_ACCESS_TOKEN", "NEW_REFRESH_TOKEN", time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
		updatedAuth, err := authRepository.UpdateAuth(context.Background(), 1, newAcessToken, newRefreshToken, &newExpiry)

		if err != nil {
			t.Fatal(err)
		} else if got, want := updatedAuth.AccessToken, newAcessToken; got != want {
			t.Fatalf("AccessToken=%v, want %v", got, want)
		} else if got, want := updatedAuth.RefreshToken, newRefreshToken; got != want {
			t.Fatalf("RefreshToken=%v, want %v", got, want)
		} else if got, want := updatedAuth.Expiry, &newExpiry; !reflect.DeepEqual(got, want) {
			t.Fatalf("Expiry=%v, want %v", got, want)
		}

		if other, err := authRepository.FindAuthByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(updatedAuth, other) {
			t.Fatalf("mismatch: %#v != %#v", updatedAuth, other)
		}

	})
}

func TestAuthRepository_FindAuth(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)
	authRepository := postgres.NewAuthRepository(db)

	user := &musichub.User{
		Name:  "jhon",
		Email: "jhon@gmai.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	auth := &musichub.Auth{
		Source:       musichub.AuthSourceGoogle,
		SourceID:     "SOURCEID",
		AccessToken:  "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
		User:         user,
		UserID:       user.ID,
	}

	if err := authRepository.CreateAuth(context.Background(), auth); err != nil {
		t.Fatal(err)
	}

	auth.User = nil

	t.Run("OKFindAuthByID", func(t *testing.T) {
		if other, err := authRepository.FindAuthByID(context.Background(), auth.ID); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(auth, other) {
			t.Fatalf("mismatch: %#v != %#v", auth, other)
		}
	})

	t.Run("OKFindAuthBySourceID", func(t *testing.T) {
		if other, err := authRepository.FindAuthBySourceID(context.Background(), auth.Source, auth.SourceID); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(auth, other) {
			t.Fatalf("mismatch: %#v != %#v", auth, other)
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		if _, err := authRepository.FindAuthByID(context.Background(), 2); musichub.ErrorCode(err) != musichub.ENOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestAuthRepository_AttachUserAuths(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	userRepository := postgres.NewUserRepository(db)

	user := &musichub.User{
		Name:  "jane",
		Email: "jane@gmail.com",
	}

	if err := userRepository.CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	t.Run("OKAttachUserAuths", func(t *testing.T) {
		authRepository := postgres.NewAuthRepository(db)

		if err := authRepository.AttachUserAuths(context.Background(), user); err != nil {
			t.Fatal(err)
		}

		if user.Auths == nil {
			t.Fatal("expected auths")
		}
	})

}

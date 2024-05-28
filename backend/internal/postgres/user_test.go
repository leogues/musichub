package postgres_test

import (
	"context"
	"reflect"
	"testing"

	musichub "github.com/leogues/MusicSyncHub"
	"github.com/leogues/MusicSyncHub/internal/postgres"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	t.Run("OKCreateUser", func(t *testing.T) {
		userRepository := postgres.NewUserRepository(db)

		user := &musichub.User{
			Name:  "jhon",
			Email: "jhon@gmail.com",
		}

		if err := userRepository.CreateUser(context.Background(), user); err != nil {
			t.Fatal(err)
		} else if got, want := user.ID, 1; got != want {
			t.Fatalf("got %v, want %v", got, want)
		} else if user.CreatedAt.IsZero() {
			t.Fatal("expected created at")
		} else if user.UpdatedAt.IsZero() {
			t.Fatal("expected updated at")
		}

		if other, err := userRepository.FindUserByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(user, other) {
			t.Fatalf("mismatch: %#v != %#v", user, other)
		}
	})

	t.Run("ErrNameRequired", func(t *testing.T) {
		userRepository := postgres.NewUserRepository(db)

		if err := userRepository.CreateUser(context.Background(), &musichub.User{}); err == nil {
			t.Fatal("expected error")
		} else if musichub.ErrorCode(err) != musichub.EINVALID || musichub.ErrorMessage(err) != "User name is required" {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestUserRepository_FindUser(t *testing.T) {
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

	t.Run("OKFindUserByID", func(t *testing.T) {
		if other, err := userRepository.FindUserByID(context.Background(), 1); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(user, other) {
			t.Fatalf("mismatch: %#v != %#v", user, other)
		}
	})

	t.Run("OKFindUserByEmail", func(t *testing.T) {
		if other, err := userRepository.FindUserByEmail(context.Background(), "jhon@gmail.com"); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(user, other) {
			t.Fatalf("mismatch: %#v != %#v", user, other)
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		userRepository := postgres.NewUserRepository(db)
		if _, err := userRepository.FindUserByID(context.Background(), 2); musichub.ErrorCode(err) != musichub.ENOTFOUND {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func TestUserRepository_AttachAuthAssociations(t *testing.T) {
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

	auth := &musichub.Auth{
		UserID: user.ID,
	}

	t.Run("OKAttachAuthAssociations", func(t *testing.T) {
		if err := userRepository.AttachAuthAssociations(context.Background(), auth); err != nil {
			t.Fatal(err)
		}

		if auth.User == nil {
			t.Fatal("expected user")
		}
	})

}

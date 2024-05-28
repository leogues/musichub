package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/leogues/MusicSyncHub/internal/postgres"
)

func TestDB(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

func MustOpenDB(tb testing.TB) *postgres.DB {
	tb.Helper()
	if err := godotenv.Load("../../.env"); err != nil {
		tb.Fatal(fmt.Errorf("ENV_FILE: %w", err))
	}

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		tb.Fatal("missing TEST_DATABASE_URL")
	}

	db := postgres.NewDB(dsn)

	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}

	return db
}

func MustCloseDB(tb testing.TB, db *postgres.DB) {
	tb.Helper()
	if err := DropAllTables(db); err != nil {
		tb.Fatal(err)
	}
	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}

func DropAllTables(db *postgres.DB) error {
	tx, err := db.BeginTx(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			DROP SCHEMA public CASCADE;
			CREATE SCHEMA public;
			GRANT ALL ON SCHEMA public TO public;
		`)

	return err
}

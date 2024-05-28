package sql

import (
	"context"
	"database/sql"
	"fmt"

	musichub "github.com/leogues/MusicSyncHub"
)

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) musichub.Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) Begin(ctx *context.Context) error {
	tx, err := t.db.BeginTx(*ctx, nil)

	if err != nil {
		return err
	}

	*ctx = musichub.NewContextWithTransaction(*ctx, tx)

	return nil
}

func (t *Transaction) Rollback(ctx context.Context) error {
	tx, ok := musichub.TransactionFromContext(ctx)

	if !ok {
		return fmt.Errorf("transaction not found")
	}

	return tx.Rollback()

}

func (t *Transaction) Commit(ctx context.Context) error {
	tx, ok := musichub.TransactionFromContext(ctx)

	if !ok {
		return fmt.Errorf("transaction not found")
	}

	return tx.Commit()
}

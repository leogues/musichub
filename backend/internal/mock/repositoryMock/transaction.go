package repositorymock

import (
	"context"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.Transaction = (*Transaction)(nil)

type Transaction struct{}

func (t *Transaction) Begin(ctx *context.Context) error {
	return nil
}

func (t *Transaction) Rollback(ctx context.Context) error {
	return nil
}

func (t *Transaction) Commit(ctx context.Context) error {
	return nil
}

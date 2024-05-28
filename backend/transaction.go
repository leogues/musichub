package musichub

import "context"

type Transaction interface {
	Begin(ctx *context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

package musichub

import (
	"context"
	"database/sql"
)

type contextProviderToken int

const (
	userContextKey          = contextProviderToken(iota + 1)
	providerTokenContextKey = contextProviderToken(iota + 1)
	txContextKey            = contextProviderToken(iota + 1)
)

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}

func UserIDFromContext(ctx context.Context) int {
	if user := UserFromContext(ctx); user != nil {
		return user.ID
	}
	return 0
}

func NewContextWithProviderToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, providerTokenContextKey, token)
}

func ProviderTokenFromContext(ctx context.Context) string {
	return ctx.Value(providerTokenContextKey).(string)
}

func NewContextWithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func TransactionFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(*sql.Tx)
	return tx, ok
}

package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*musichub.User, error) {
	users, _, err := r.FindUsers(ctx, musichub.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "User not found."}
	}
	return users[0], nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*musichub.User, error) {
	a, _, err := r.FindUsers(ctx, musichub.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *musichub.User) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt

	if err := user.Validate(); err != nil {
		return err
	}

	var email *string
	if user.Email != "" {
		email = &user.Email
	}

	apiKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, apiKey); err != nil {
		return err
	}
	user.APIKey = hex.EncodeToString(apiKey)

	result := tx.QueryRowContext(ctx, `
		INSERT INTO public.users (
			name,
			email,
			api_key,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		user.Name,
		email,
		user.APIKey,
		(*NullTime)(&user.CreatedAt),
		(*NullTime)(&user.UpdatedAt),
	)
	err = result.Err()
	if err != nil {
		return err
	}

	var id int
	err = result.Scan(&id)
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

func (r *UserRepository) FindUsers(ctx context.Context, filter musichub.UserFilter) (_ []*musichub.User, n int, err error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	index := 1
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}
	if v := filter.Email; v != nil {
		where, args = append(where, "email = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}

	if v := filter.APIKey; v != nil {
		where, args = append(where, "api_key = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    name,
		    email,
				api_key,
		    created_at,
		    updated_at,
				COUNT(*) OVER()
		FROM public.users
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*musichub.User, 0)
	for rows.Next() {
		var email sql.NullString
		var user musichub.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&email,
			&user.APIKey,
			(*NullTime)(&user.CreatedAt),
			(*NullTime)(&user.UpdatedAt),
			&n,
		); err != nil {
			return nil, 0, err
		}

		if email.Valid {
			user.Email = email.String
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

func (r *UserRepository) AttachAuthAssociations(ctx context.Context, auth *musichub.Auth) (err error) {
	if auth.User, err = r.FindUserByID(ctx, auth.UserID); err != nil {
		return fmt.Errorf("attach auth users: %w", err)
	}
	return nil
}

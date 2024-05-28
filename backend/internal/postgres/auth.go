package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
)

var _ musichub.AuthRepository = (*AuthRepository)(nil)

type AuthRepository struct {
	db *DB
}

func NewAuthRepository(db *DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindAuthByID(ctx context.Context, id int) (*musichub.Auth, error) {
	auths, err := r.findAuths(ctx, musichub.AuthFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(auths) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "Auth not found."}
	}
	return auths[0], nil
}

func (r *AuthRepository) FindAuthBySourceID(ctx context.Context, source string, sourceID string) (*musichub.Auth, error) {
	auths, err := r.findAuths(ctx, musichub.AuthFilter{Source: &source, SourceID: &sourceID})
	if err != nil {
		return nil, err
	} else if len(auths) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "Auth not found."}
	}

	return auths[0], nil
}

func (r *AuthRepository) UpdateAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.Auth, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := r.FindAuthByID(ctx, id)
	if err != nil {
		return auth, err
	}

	auth.AccessToken = accessToken
	auth.RefreshToken = refreshToken
	auth.Expiry = expiry
	auth.UpdatedAt = tx.now

	if err := auth.Validate(); err != nil {
		return auth, err
	}

	var expiryStr *string
	if auth.Expiry != nil {
		v := auth.Expiry.Format(time.RFC3339)
		expiryStr = &v
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE public.auths
		SET access_token = $1,
		    refresh_token = $2,
		    expiry = $3,
		    updated_at = $4
		WHERE id = $5
	`,
		auth.AccessToken,
		auth.RefreshToken,
		expiryStr,
		(*NullTime)(&auth.UpdatedAt),
		id,
	); err != nil {
		return auth, err
	}

	return auth, nil
}

func (r *AuthRepository) CreateAuth(ctx context.Context, auth *musichub.Auth) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	auth.CreatedAt = tx.now
	auth.UpdatedAt = auth.CreatedAt

	if err := auth.Validate(); err != nil {
		return err
	}

	var expiry *string
	if auth.Expiry != nil {
		tmp := auth.Expiry.Format(time.RFC3339)
		expiry = &tmp
	}

	result := tx.QueryRowContext(ctx, `
		INSERT INTO public.auths (
			user_id,
			source,
			source_id,
			access_token,
			refresh_token,
			expiry,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		auth.UserID,
		auth.Source,
		auth.SourceID,
		auth.AccessToken,
		auth.RefreshToken,
		expiry,
		(*NullTime)(&auth.CreatedAt),
		(*NullTime)(&auth.UpdatedAt),
	)
	err = result.Err()
	if err != nil {
		return err
	}

	var id int
	if err = result.Scan(&id); err != nil {
		return err
	}
	auth.ID = int(id)

	return nil
}

func (r *AuthRepository) findAuths(ctx context.Context, filter musichub.AuthFilter) ([]*musichub.Auth, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	index := 1
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}
	if v := filter.UserID; v != nil {
		where, args = append(where, "user_id = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}
	if v := filter.SourceID; v != nil {
		where, args = append(where, "source_id = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}
	if v := filter.Source; v != nil {
		where, args = append(where, "source = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
				id,
				user_id,
				source,
				source_id,
				access_token,
				refresh_token,
				expiry,
				created_at,
				updated_at
		FROM public.auths
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	auths := make([]*musichub.Auth, 0)
	for rows.Next() {
		var auth musichub.Auth
		var expiry sql.NullString
		if err = rows.Scan(
			&auth.ID,
			&auth.UserID,
			&auth.Source,
			&auth.SourceID,
			&auth.AccessToken,
			&auth.RefreshToken,
			&expiry,
			(*NullTime)(&auth.CreatedAt),
			(*NullTime)(&auth.UpdatedAt),
		); err != nil {
			return nil, err
		}

		if expiry.Valid {
			if t, _ := time.Parse(time.RFC3339, expiry.String); !t.IsZero() {
				auth.Expiry = &t
			}
		}

		auths = append(auths, &auth)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return auths, nil
}

func (r *AuthRepository) AttachUserAuths(ctx context.Context, user *musichub.User) (err error) {
	if user.Auths, err = r.findAuths(ctx, musichub.AuthFilter{UserID: &user.ID}); err != nil {
		return fmt.Errorf("attach user auth: %w", err)
	}
	return nil
}

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

var _ musichub.ProviderAuthRepository = (*ProviderAuthRepository)(nil)

type ProviderAuthRepository struct {
	db *DB
}

func NewProviderAuthRepository(db *DB) *ProviderAuthRepository {
	return &ProviderAuthRepository{db: db}
}

func (r *ProviderAuthRepository) FindProviderAuthByID(ctx context.Context, id int) (*musichub.ProviderAuth, error) {
	providerAuths, err := r.findProviderAuths(ctx, musichub.ProviderAuthFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(providerAuths) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "ProviderAuth not found."}
	}

	return providerAuths[0], nil
}

func (r *ProviderAuthRepository) FindProviderAuthBySource(ctx context.Context, userId int, source string) (*musichub.ProviderAuth, error) {
	providerAuths, err := r.findProviderAuths(ctx, musichub.ProviderAuthFilter{UserID: &userId, Source: &source})
	if err != nil {
		return nil, err
	} else if len(providerAuths) == 0 {
		return nil, &musichub.Error{Code: musichub.ENOTFOUND, Message: "ProviderAuth not found."}
	}

	return providerAuths[0], nil
}

func (r *ProviderAuthRepository) UpdateProviderAuth(ctx context.Context, id int, accessToken, refreshToken string, expiry *time.Time) (*musichub.ProviderAuth, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	providerAuth, err := r.FindProviderAuthByID(ctx, id)
	if err != nil {
		return providerAuth, err
	}

	providerAuth.AccessToken = accessToken
	providerAuth.RefreshToken = refreshToken
	providerAuth.Expiry = expiry
	providerAuth.UpdatedAt = tx.now

	if err := providerAuth.Validate(); err != nil {
		return providerAuth, err
	}

	var expiryStr *string
	if providerAuth.Expiry != nil {
		v := providerAuth.Expiry.Format(time.RFC3339)
		expiryStr = &v
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE public.provider_auths
		SET access_token = $1,
		    refresh_token = $2,
		    expiry = $3,
		    updated_at = $4
		WHERE id = $5
	`,
		providerAuth.AccessToken,
		providerAuth.RefreshToken,
		expiryStr,
		(*NullTime)(&providerAuth.UpdatedAt),
		id,
	); err != nil {
		return providerAuth, err
	}

	return providerAuth, nil
}

func (r *ProviderAuthRepository) CreateProviderAuth(ctx context.Context, auth *musichub.ProviderAuth) error {
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
		INSERT INTO public.provider_auths (
			user_id,
			source,
			access_token,
			refresh_token,
			expiry,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		auth.UserID,
		auth.Source,
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
	if err := result.Scan(&id); err != nil {
		return err
	}

	auth.ID = id

	return nil
}

func (r *ProviderAuthRepository) DeleteProviderAuth(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	if providerAuth, err := r.FindProviderAuthByID(ctx, id); err != nil {
		return err
	} else if providerAuth.UserID != musichub.UserIDFromContext(ctx) {
		return musichub.Errorf(musichub.EUNAUTHORIZED, "You are not allowed to delete this provider auth.")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM public.provider_auths WHERE id = $1`, id); err != nil {
		return err
	}

	return nil

}

func (r *ProviderAuthRepository) findProviderAuths(ctx context.Context, filter musichub.ProviderAuthFilter) ([]*musichub.ProviderAuth, error) {
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

	if v := filter.Source; v != nil {
		where, args = append(where, "source = $"+strconv.Itoa(index)), append(args, *v)
		index++
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
				id,
				user_id,
				source,
				access_token,
				refresh_token,
				expiry,
				created_at,
				updated_at
		FROM public.provider_auths
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providerAuths := make([]*musichub.ProviderAuth, 0)
	for rows.Next() {
		var providerAuth musichub.ProviderAuth
		var expiry sql.NullString
		if err = rows.Scan(
			&providerAuth.ID,
			&providerAuth.UserID,
			&providerAuth.Source,
			&providerAuth.AccessToken,
			&providerAuth.RefreshToken,
			&expiry,
			(*NullTime)(&providerAuth.CreatedAt),
			(*NullTime)(&providerAuth.UpdatedAt),
		); err != nil {
			return nil, err
		}

		if expiry.Valid {
			if t, _ := time.Parse(time.RFC3339, expiry.String); !t.IsZero() {
				providerAuth.Expiry = &t
			}
		}

		providerAuths = append(providerAuths, &providerAuth)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return providerAuths, nil
}

func (r *ProviderAuthRepository) AttachUserProviderAuths(ctx context.Context, user *musichub.User) (err error) {
	if user.ProviderAuths, err = r.findProviderAuths(ctx, musichub.ProviderAuthFilter{UserID: &user.ID}); err != nil {
		return fmt.Errorf("attach user auth: %w", err)
	}

	return nil

}

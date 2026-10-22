package auth

import (
	"context"
	"time"

	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(
	ctx context.Context,
	userID uuid.UUID,
	identityID uuid.UUID,
	purpose string,
	target string,
	tokenHash string,
	expiresAt time.Time,
) (*authModel.AuthToken, error) {
	var token authModel.AuthToken

	err := r.db.GetContext(ctx, &token, `
        INSERT INTO auth_tokens (
            user_id, identity_id, purpose, target, token_hash, expires_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING
            id, user_id, identity_id, purpose, target, token_hash,
            expires_at, consumed_at, created_at, updated_at, deleted_at
    `, userID, identityID, purpose, target, tokenHash, expiresAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*authModel.AuthToken, error) {
	var token authModel.AuthToken

	err := r.db.GetContext(ctx, &token, `
        SELECT
            id, user_id, identity_id, purpose, target, token_hash,
            expires_at, consumed_at, created_at, updated_at, deleted_at
        FROM auth_tokens
        WHERE token_hash = $1
          AND deleted_at IS NULL
    `, tokenHash)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) GetActiveByTargetPurpose(
	ctx context.Context,
	target string,
	purpose string,
) (*authModel.AuthToken, error) {
	var token authModel.AuthToken

	err := r.db.GetContext(ctx, &token, `
        SELECT
            id, user_id, identity_id, purpose, target, token_hash,
            expires_at, consumed_at, created_at, updated_at, deleted_at
        FROM auth_tokens
        WHERE target = $1
          AND purpose = $2
          AND consumed_at IS NULL
          AND deleted_at IS NULL
          AND expires_at > NOW()
        ORDER BY created_at DESC
        LIMIT 1
    `, target, purpose)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) Consume(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE auth_tokens
        SET consumed_at = NOW()
        WHERE id = $1
          AND consumed_at IS NULL
          AND deleted_at IS NULL
          AND expires_at > NOW()
    `, id)

	return err
}

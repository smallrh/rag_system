package auth

import (
	"context"
	"encoding/json"
	"time"

	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	userID uuid.UUID,
	identityID uuid.UUID,
	provider string,
	refreshTokenHash string,
	deviceInfo json.RawMessage,
	expiresAt time.Time,
) (*authModel.AuthSession, error) {
	var session authModel.AuthSession

	if len(deviceInfo) == 0 {
		deviceInfo = json.RawMessage(`{}`)
	}

	err := r.db.GetContext(ctx, &session, `
        INSERT INTO auth_sessions (
            user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING
            id, user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at,
            last_seen_at, rotated_at, replaced_by_session_id, revoked_at,
            created_at, updated_at, deleted_at
    `, userID, identityID, provider, refreshTokenHash, deviceInfo, expiresAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.AuthSession, error) {
	var session authModel.AuthSession

	err := r.db.GetContext(ctx, &session, `
        SELECT
            id, user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at,
            last_seen_at, rotated_at, replaced_by_session_id, revoked_at,
            created_at, updated_at, deleted_at
        FROM auth_sessions
        WHERE id = $1
          AND deleted_at IS NULL
    `, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetByRefreshTokenHash(
	ctx context.Context,
	refreshTokenHash string,
) (*authModel.AuthSession, error) {
	var session authModel.AuthSession

	err := r.db.GetContext(ctx, &session, `
        SELECT
            id, user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at,
            last_seen_at, rotated_at, replaced_by_session_id, revoked_at,
            created_at, updated_at, deleted_at
        FROM auth_sessions
        WHERE refresh_token_hash = $1
          AND deleted_at IS NULL
    `, refreshTokenHash)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetActiveByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.AuthSession, error) {
	var session authModel.AuthSession

	err := r.db.GetContext(ctx, &session, `
        SELECT
            id, user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at,
            last_seen_at, rotated_at, replaced_by_session_id, revoked_at,
            created_at, updated_at, deleted_at
        FROM auth_sessions
        WHERE id = $1
          AND revoked_at IS NULL
          AND deleted_at IS NULL
          AND expires_at > NOW()
        FOR UPDATE
    `, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) ListActiveByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]authModel.AuthSession, error) {
	var sessions []authModel.AuthSession

	err := r.db.SelectContext(ctx, &sessions, `
        SELECT
            id, user_id, identity_id, provider,
            refresh_token_hash, device_info, expires_at,
            last_seen_at, rotated_at, replaced_by_session_id, revoked_at,
            created_at, updated_at, deleted_at
        FROM auth_sessions
        WHERE user_id = $1
          AND revoked_at IS NULL
          AND deleted_at IS NULL
          AND expires_at > NOW()
        ORDER BY id
    `, userID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) TouchLastSeen(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE auth_sessions
        SET last_seen_at = NOW()
        WHERE id = $1
          AND deleted_at IS NULL
    `, id)

	return err
}

func (r *SessionRepository) Rotate(
	ctx context.Context,
	tx *sqlx.Tx,
	id uuid.UUID,
	newRefreshTokenHash string,
	replacedBySessionID uuid.UUID,
) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE auth_sessions
        SET refresh_token_hash = $2,
            rotated_at = NOW(),
            replaced_by_session_id = $3
        WHERE id = $1
          AND revoked_at IS NULL
          AND deleted_at IS NULL
    `, id, newRefreshTokenHash, replacedBySessionID)

	return err
}

func (r *SessionRepository) Revoke(
	ctx context.Context,
	tx *sqlx.Tx,
	id uuid.UUID,
) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE auth_sessions
        SET revoked_at = NOW()
        WHERE id = $1
          AND revoked_at IS NULL
          AND deleted_at IS NULL
    `, id)

	return err
}

func (r *SessionRepository) RevokeAllByUserID(
	ctx context.Context,
	tx *sqlx.Tx,
	userID uuid.UUID,
) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE auth_sessions
        SET revoked_at = NOW()
        WHERE user_id = $1
          AND revoked_at IS NULL
          AND deleted_at IS NULL
    `, userID)

	return err
}

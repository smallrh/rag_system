package auth

import (
	"context"
	"encoding/json"
	"time"

	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IdentityRepository struct {
	db *sqlx.DB
}

func NewIdentityRepository(db *sqlx.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

func (r *IdentityRepository) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	userID uuid.UUID,
	provider string,
	identifier string,
	metadata json.RawMessage,
) (*authModel.AuthIdentity, error) {
	var identity authModel.AuthIdentity

	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	err := tx.GetContext(ctx, &identity, `
        INSERT INTO auth_identities (user_id, provider, identifier, metadata)
        VALUES ($1, $2, $3, $4)
        RETURNING
            id, user_id, provider, identifier,
            verified_at, status, metadata,
            created_at, updated_at, deleted_at
    `, userID, provider, identifier, metadata)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (r *IdentityRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.AuthIdentity, error) {
	var identity authModel.AuthIdentity

	err := r.db.GetContext(ctx, &identity, `
        SELECT
            id, user_id, provider, identifier,
            verified_at, status, metadata,
            created_at, updated_at, deleted_at
        FROM auth_identities
        WHERE id = $1
          AND deleted_at IS NULL
    `, id)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (r *IdentityRepository) GetByProviderIdentifier(
	ctx context.Context,
	provider string,
	identifier string,
) (*authModel.AuthIdentity, error) {
	var identity authModel.AuthIdentity

	err := r.db.GetContext(ctx, &identity, `
        SELECT
            id, user_id, provider, identifier,
            verified_at, status, metadata,
            created_at, updated_at, deleted_at
        FROM auth_identities
        WHERE provider = $1
          AND LOWER(identifier) = LOWER($2)
          AND deleted_at IS NULL
    `, provider, identifier)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (r *IdentityRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]authModel.AuthIdentity, error) {
	var identities []authModel.AuthIdentity

	err := r.db.SelectContext(ctx, &identities, `
        SELECT
            id, user_id, provider, identifier,
            verified_at, status, metadata,
            created_at, updated_at, deleted_at
        FROM auth_identities
        WHERE user_id = $1
          AND deleted_at IS NULL
        ORDER BY id
    `, userID)
	if err != nil {
		return nil, err
	}

	return identities, nil
}

func (r *IdentityRepository) MarkVerified(
	ctx context.Context,
	tx *sqlx.Tx,
	id uuid.UUID,
	verifiedAt time.Time,
) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE auth_identities
        SET verified_at = $2
        WHERE id = $1
          AND deleted_at IS NULL
    `, id, verifiedAt)

	return err
}

func (r *IdentityRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE auth_identities
        SET status = $2
        WHERE id = $1
          AND deleted_at IS NULL
    `, id, status)

	return err
}

package auth

import (
	"context"
	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CredentialRepository struct {
	db *sqlx.DB
}

func NewCredentialRepository(db *sqlx.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

func (r *CredentialRepository) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	userID uuid.UUID,
	identityID *uuid.UUID,
	credentialType string,
	secretHash *string,
) (*authModel.AuthCredential, error) {
	var credential authModel.AuthCredential

	err := tx.GetContext(ctx, &credential, `
        INSERT INTO auth_credentials (user_id, identity_id, credential_type, secret_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING
            id, user_id, identity_id, credential_type, secret_hash, metadata,
            created_at, updated_at, deleted_at
    `, userID, identityID, credentialType, secretHash)
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (r *CredentialRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.AuthCredential, error) {
	var credential authModel.AuthCredential

	err := r.db.GetContext(ctx, &credential, `
        SELECT
            id, user_id, identity_id, credential_type, secret_hash, metadata,
            created_at, updated_at, deleted_at
        FROM auth_credentials
        WHERE id = $1
          AND deleted_at IS NULL
    `, id)
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (r *CredentialRepository) GetPasswordByUserIDAndCredentialType(
	ctx context.Context,
	userID uuid.UUID,
	credentialType string,
) (*authModel.AuthCredential, error) {
	var credential authModel.AuthCredential

	err := r.db.GetContext(ctx, &credential, `
        SELECT
            id, user_id, identity_id, credential_type, secret_hash, metadata,
            created_at, updated_at, deleted_at
        FROM auth_credentials
        WHERE user_id = $1
          AND credential_type = $2
          AND deleted_at IS NULL
    `, userID, credentialType)
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (r *CredentialRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]authModel.AuthCredential, error) {
	var credentials []authModel.AuthCredential

	err := r.db.SelectContext(ctx, &credentials, `
        SELECT
            id, user_id, identity_id, credential_type, secret_hash, metadata,
            created_at, updated_at, deleted_at
        FROM auth_credentials
        WHERE user_id = $1
          AND deleted_at IS NULL
        ORDER BY id
    `, userID)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (r *CredentialRepository) UpdateSecretHash(
	ctx context.Context,
	tx *sqlx.Tx,
	id uuid.UUID,
	secretHash string,
) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE auth_credentials
        SET secret_hash = $2
        WHERE id = $1
          AND deleted_at IS NULL
    `, id, secretHash)

	return err
}

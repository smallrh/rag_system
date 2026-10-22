package auth

import (
	"context"

	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	displayName string,
) (*authModel.User, error) {
	var user authModel.User

	err := r.db.GetContext(ctx, &user, `
        INSERT INTO users (display_name)
        VALUES (NULLIF($1, ''))
        RETURNING
            id, display_name, status,
            created_at, updated_at, deleted_at
    `, displayName)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.User, error) {
	var user authModel.User

	err := r.db.GetContext(ctx, &user, `
        SELECT
            id, display_name, status,
            created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1
          AND deleted_at IS NULL
    `, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetActiveByID(
	ctx context.Context,
	id uuid.UUID,
) (*authModel.User, error) {
	var user authModel.User

	err := r.db.GetContext(ctx, &user, `
        SELECT
            id, display_name, status,
            created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1
          AND status = 'active'
          AND deleted_at IS NULL
        FOR UPDATE
    `, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE users
        SET status = $2,
            updated_at = NOW()
        WHERE id = $1
          AND deleted_at IS NULL
    `, id, status)

	return err
}

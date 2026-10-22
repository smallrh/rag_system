package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	IdentityID uuid.UUID `db:"identity_id"`
	Purpose    string    `db:"purpose"`
	Target     string    `db:"target"`
	TokenHash  string    `db:"token_hash"`
	ExpiresAt  time.Time `db:"expires_at"`
	ConsumedAt time.Time `db:"consumed_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	DeletedAt  time.Time `db:"deleted_at"`
}

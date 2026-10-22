package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthIdentity struct {
	ID         uuid.UUID       `db:"id"`
	UserID     uuid.UUID       `db:"user_id"`
	Provider   string          `db:"provider"`
	Identifier string          `db:"identifier"`
	VerifiedAt time.Time       `db:"verified_at"`
	Status     string          `db:"status"`
	Metadata   json.RawMessage `db:"metadata"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	DeletedAt  time.Time       `db:"deleted_at"`
}

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
	VerifiedAt *time.Time      `db:"verified_at"`
	Status     string          `db:"status"`
	Metadata   json.RawMessage `db:"metadata"`
}

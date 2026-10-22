package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthCredential struct {
	ID             uuid.UUID       `db:"id"`
	UserID         uuid.UUID       `db:"user_id"`
	IdentityID     uuid.UUID       `db:"identity_id"`
	CredentialType string          `db:"credential_type"`
	SecretHash     string          `db:"secret_hash"`
	Metadata       json.RawMessage `db:"metadata"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	DeletedAt      time.Time       `db:"deleted_at"`
}

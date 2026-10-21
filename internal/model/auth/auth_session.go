package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID                  uuid.UUID       `db:"id"`
	UserID              uuid.UUID       `db:"user_id"`
	IdentityID          *uuid.UUID      `db:"identity_id"`
	Provider            *string         `db:"provider"`
	RefreshTokenHash    string          `db:"refresh_token_hash"`
	DeviceInfo          json.RawMessage `db:"device_info"`
	ExpiresAt           time.Time       `db:"expires_at"`
	LastSeenAt          *time.Time      `db:"last_seen_at"`
	RotatedAt           *time.Time      `db:"rotated_at"`
	ReplacedBySessionID *uuid.UUID      `db:"replaced_by_session_id"`
	RevokedAt           *time.Time      `db:"revoked_at"`
}

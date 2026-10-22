package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthAuditLog struct {
	ID            int64     `db:"id"`
	UserID        uuid.UUID `db:"user_id"`
	IdentityID    uuid.UUID `db:"identity_id"`
	EventType     string    `db:"event_type"`
	Success       bool      `db:"success"`
	Identifier    string    `db:"identifier"`
	IPAddress     string    `db:"ip_address"`
	UserAgent     string    `db:"user_agent"`
	FailureReason string    `db:"failure_reason"`
	CreatedAt     time.Time `db:"created_at"`
}

package auth

import "github.com/google/uuid"

type AuthCredential struct {
	ID             uuid.UUID  `db:"id"`
	UserID         uuid.UUID  `db:"user_id"`
	IdentityID     *uuid.UUID `db:"identity_id"`
	CredentialType string     `db:"credential_type"`
	SecretHash     *string    `db:"secret_hash"`
}

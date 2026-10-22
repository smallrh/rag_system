package login

import (
	"context"

	"github.com/google/uuid"

	passwordsvc "go_study/internal/services/auth/password"
)

func (s *Service) checkPassword(ctx context.Context, userID uuid.UUID, plain string) error {
	if plain == "" || len(plain) > 72 {
		return ErrInvalidCredentials
	}
	cred, err := s.repo.Credential.GetPasswordByUserIDAndCredentialType(ctx, userID, "password")
	if err != nil || cred == nil {
		return ErrInvalidCredentials
	}
	if isUnset(cred.SecretHash) {
		return ErrInvalidCredentials
	}
	hash := derefString(cred.SecretHash)
	if hash == "" {
		return ErrInvalidCredentials
	}
	return passwordsvc.VerifyPassword(plain, hash)
}

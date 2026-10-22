package login

import (
	"context"
	"crypto/subtle"
	"strings"

	authModel "go_study/internal/model/auth"
)

func (s *Service) consumeAuthCode(
	ctx context.Context,
	identity *authModel.AuthIdentity,
	code string,
) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return ErrInvalidAuthCode
	}

	var purpose string
	switch identity.Provider {
	case "email":
		purpose = "email_login"
	case "phone":
		purpose = "phone_login"
	default:
		return ErrUnsupportedType
	}

	target := strings.ToLower(strings.TrimSpace(identity.Identifier))
	tok, err := s.repo.Token.GetActiveByTargetPurpose(ctx, target, purpose)
	if err != nil || tok == nil {
		return ErrInvalidAuthCode
	}

	if subtle.ConstantTimeCompare([]byte(tok.TokenHash), []byte(sha256Hex(code))) != 1 {
		return ErrInvalidAuthCode
	}

	// Consume 签名是 (ctx, id)；请确保 repo 实现是原子的：
	// UPDATE auth_tokens SET consumed_at = NOW()
	//  WHERE id = $1 AND consumed_at IS NULL AND deleted_at IS NULL AND expires_at > NOW()
	return s.repo.Token.Consume(ctx, tok.ID)
}

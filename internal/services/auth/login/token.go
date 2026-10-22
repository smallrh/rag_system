package login

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	authModel "go_study/internal/model/auth"
)

func (s *Service) issueAccessToken(
	user *authModel.User,
	identity *authModel.AuthIdentity,
	now, expiresAt time.Time,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"ident_id": identity.ID.String(),
		"provider": identity.Provider,
		"iat":      now.Unix(),
		"exp":      expiresAt.Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

package auth

import "context"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"` // email或者phone
	AuthCode string `json:"auth_code"`
}

func (s *AuthService) Register(
	ctx context.Context,
	req RegisterRequest,
) error {
	return nil
}

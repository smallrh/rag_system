package meta

import "go_study/internal/services/auth"

type LoginRequestMeta struct {
	auth.LoginRequest
	auth.VerifyCodeRequest
}

type RegisterRequestMeta struct {
	auth.RegisterRequest
	auth.VerifyCodeRequest
}

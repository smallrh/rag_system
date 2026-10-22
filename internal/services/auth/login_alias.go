package auth

import "go_study/internal/services/auth/login"

// 保留原有对外类型和错误，实际实现位于 services/auth/login 子包。
type (
	LoginRequest  = login.LoginRequest
	LoginResponse = login.LoginResponse
)

const (
	LoginTypePassword = login.LoginTypePassword
	LoginTypeCode     = login.LoginTypeCode

	RefreshTokenTTL = login.RefreshTokenTTL
	AccessTokenTTL  = login.AccessTokenTTL
)

var (
	ErrInvalidCredentials = login.ErrInvalidCredentials
	ErrInvalidAuthCode    = login.ErrInvalidAuthCode
	ErrIdentityInactive   = login.ErrIdentityInactive
	ErrUserNotActive      = login.ErrUserNotActive
	ErrUnsupportedType    = login.ErrUnsupportedType
)

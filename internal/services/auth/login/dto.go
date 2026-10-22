package login

import (
	"errors"
	"time"
)

type LoginRequest struct {
	Provider   string         `json:"provider"`
	Username   string         `json:"username"`
	Password   string         `json:"password,omitempty"`
	Type       string         `json:"type"` // password | code
	AuthCode   string         `json:"auth_code,omitempty"`
	DeviceInfo map[string]any `json:"device_info,omitempty"`
	IPAddress  string         `json:"ip_address,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuthCode    = errors.New("invalid or expired auth code")
	ErrIdentityInactive   = errors.New("identity disabled or deleted")
	ErrUserNotActive      = errors.New("user not active")
	ErrUnsupportedType    = errors.New("unsupported login type")
)

const (
	LoginTypePassword = "password"
	LoginTypeCode     = "code"

	RefreshTokenTTL = 30 * 24 * time.Hour
	AccessTokenTTL  = 2 * time.Hour
	refreshTokenN   = 32 // 256-bit 随机
)

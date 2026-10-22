package auth

import "go_study/internal/services/auth/password"

// 保留原有对外函数，实际实现位于 services/auth/password 子包。
var ErrInvalidPassword = password.ErrInvalidPassword

func HashPassword(plain string) (string, error) {
	return password.HashPassword(plain)
}

func VerifyPassword(plain, hashedPassword string) error {
	return password.VerifyPassword(plain, hashedPassword)
}

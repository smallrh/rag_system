package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("invalid password")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidPassword
	}

	// 哈希格式错误等问题属于凭据数据异常，
	// 不应伪装成普通的密码错误。
	return err
}

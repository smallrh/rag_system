package auth

import (
	repoAuth "go_study/internal/repository/auth"
	"go_study/internal/services/auth/login"
)

type AuthService struct {
	repo *repoAuth.Repository
	*login.Service
}

func NewAuthService(repo *repoAuth.Repository, jwtSecret []byte) *AuthService {
	return &AuthService{
		repo:    repo,
		Service: login.NewService(repo, jwtSecret),
	}
}

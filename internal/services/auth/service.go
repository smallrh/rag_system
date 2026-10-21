package auth

import repoAuth "go_study/internal/repository/auth"

type AuthService struct {
	repo *repoAuth.Repository
}

func NewAuthService(repo *repoAuth.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

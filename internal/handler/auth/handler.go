package auth

import serviceAuth "go_study/internal/services/auth"

type AuthHandler struct {
	service *serviceAuth.AuthService
}

func NewAuthHandler(service *serviceAuth.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

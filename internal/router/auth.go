package router

import (
	handlerAuth "go_study/internal/handler/auth"

	"github.com/labstack/echo/v5"
)

func RegisterAuthRouter(
	e *echo.Echo,
	h *handlerAuth.AuthHandler,
) {
	g := e.Group("/auth")

	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
	g.GET("/captcha", h.GetCaptcha)
	g.POST("/verify_captcha", h.VerifyCaptcha)
}

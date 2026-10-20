package router

import (
	handlerAuth "go_study/internal/handler/auth"

	"github.com/labstack/echo/v5"
)

func RegisterAuthRouter(e *echo.Echo) {
	g := e.Group("/auth")
	g.POST("/login", handlerAuth.Login)
	g.POST("/register", handlerAuth.Register)
}

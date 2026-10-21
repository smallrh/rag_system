package router

import (
	handlerAuth "go_study/internal/handler/auth"

	"github.com/labstack/echo/v5"
)

func RegisterRouters(e *echo.Echo, h *handlerAuth.AuthHandler) {
	RegisterAuthRouter(e, h)
	// RegisterUserRouter(e)
	// RegisterOrderRouter(e)
}

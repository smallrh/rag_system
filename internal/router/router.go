package router

import "github.com/labstack/echo/v5"

func RegisterRouters(e *echo.Echo) {
	RegisterAuthRouter(e)
	// RegisterUserRouter(e)
	// RegisterOrderRouter(e)
}

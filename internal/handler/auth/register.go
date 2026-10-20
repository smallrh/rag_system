package auth

import (
	"go_study/internal/handler/meta"
	"go_study/internal/services/auth"
	"go_study/internal/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

func Register(c *echo.Context) error {
	data, err := utils.Convert[meta.RegisterRequestMeta](c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}
	success, err := auth.VerifyCode(data.VerifyCodeRequest)
	if err != nil {
		return err
	}
	if success {
		auth.Register(data.RegisterRequest)
	}
	if !success {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid verification code",
		)
	}

	return nil
}

package auth

import (
	"go_study/internal/handler/meta"
	"go_study/internal/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *AuthHandler) Register(c *echo.Context) error {
	data, err := utils.Convert[meta.RegisterRequestMeta](c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}
	success := h.service.VerifyCode(data.VerifyCodeRequest)

	if success {
		err := h.service.Register(c.Request().Context(), data.RegisterRequest)
		if err != nil {
			return err
		}
	}
	if !success {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid verification code",
		)
	}
	return nil
}

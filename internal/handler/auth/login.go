package auth

import (
	"go_study/internal/handler/meta"
	"go_study/internal/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

var (
	ErrorBlock = echo.NewHTTPError(
		http.StatusBadRequest,
		"invalid verification",
	)
)

func (h *AuthHandler) Login(c *echo.Context) error {
	req, err := utils.Convert[meta.LoginRequestMeta](c.Request().Body)
	if err != nil {
		return ErrorBlock
	}
	switch req.Type {
	case "password":
		if !h.service.VerifyCode(req.VerifyCodeRequest) {
			return ErrorBlock
		}
	case "code":
		// 验证码登录，无需额外校验图形验证码
	default:
		return ErrorBlock
	}
	data, _ := h.service.Login(c.Request().Context(), req.LoginRequest)
	return c.JSON(http.StatusOK, data)
}

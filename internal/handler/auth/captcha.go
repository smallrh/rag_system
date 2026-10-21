package auth

import (
	"go_study/internal/utils"
	"net/http"

	captchaService "go_study/internal/services/auth"

	"github.com/labstack/echo/v5"
)

func (h *AuthHandler) GetCaptcha(c *echo.Context) error {
	data, captchaID, err := h.service.GetCaptcha()
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate captcha",
		)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"captcha_id": captchaID,
		"image":      data,
	})
}

func (h *AuthHandler) VerifyCaptcha(c *echo.Context) error {
	req, _ := utils.Convert[captchaService.VerifyCaptchaRequest](c.Request().Body)
	data := h.service.VerifyCode(req)
	return c.JSON(http.StatusOK, map[string]any{
		"captcha_id": req.CaptchaId,
		"success":    data,
	})
}

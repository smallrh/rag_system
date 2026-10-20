package auth

import (
	"go_study/internal/utils"
	"net/http"

	"github.com/labstack/echo/v5"
	captchaService "go_study/internal/services/auth"
)

func GetCaptcha(c *echo.Context) error {
	data, captchaID, err := captchaService.GetCaptcha()
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

func VerifyCode(c *echo.Context) error {
	req, _ := utils.Convert[captchaService.VerifyCodeRequest](c.Request().Body)
	data := captchaService.VerifyCode(req)
	return c.JSON(http.StatusOK, map[string]any{
		"captcha_id": req.CaptchaId,
		"success":    data,
	})
}

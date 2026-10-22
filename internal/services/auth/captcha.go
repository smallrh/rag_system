package auth

import (
	"bytes"
	"encoding/base64"

	"github.com/dchest/captcha"
)

func (s *AuthService) GetCaptcha() (string, string, error) {
	id := captcha.New()
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, id, 240, 80); err != nil {
		return "", "", err
	}
	image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return "data:image/png;base64," + image, id, nil
}

type VerifyCodeRequest struct {
	CaptchaId   string `json:"captcha_id"`
	CaptchaCode string `json:"Captcha_code"`
}

func (s *AuthService) VerifyCode(req VerifyCodeRequest) bool {
	data := captcha.VerifyString(req.CaptchaId, req.CaptchaCode)
	return data
}

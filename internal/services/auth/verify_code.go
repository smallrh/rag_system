package auth

type VerifyCodeRequest struct {
	Type     string `json:"type"`
	UserName string `json:"user_name"`
	Code     string `json:"code"`
}

func VerifyCode(req VerifyCodeRequest) (bool, error) {
	return true, nil
}

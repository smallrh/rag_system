package auth

import (
	"context"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"` // code 或者 password
	AuthCode string `json:"auth_code"`
}
type LoginResponse struct {
}

func (s *AuthService) Login(
	ctx context.Context,
	req LoginRequest,
) (LoginResponse, error) {
	// 1. 根据登录方式确定身份标识
	// 2. 通过 s.users 查询用户及认证凭据
	// 3. 校验密码或验证码
	// 4. 根据需要签发 Token
	// 5. 返回登录结果

	return LoginResponse{}, nil
}

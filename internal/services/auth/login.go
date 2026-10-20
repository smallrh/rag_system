package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(req LoginRequest) map[string]error {
	return nil
}

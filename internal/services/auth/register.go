package auth

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(req RegisterRequest) map[string]error {
	return nil
}

package auth

type LoginRequestDto struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginUserDto struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginResponseDto struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresAt   string       `json:"expires_at"`
	User        LoginUserDto `json:"user"`
}

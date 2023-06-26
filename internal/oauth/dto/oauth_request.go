package oauth

type LoginRequestBody struct {
	Email        string `josn:"email" binding:"email"`
	Password     string `josn:"password" binding:"required"`
	ClientID     string `josn:"client_id" binding:"required"`
	ClientSecret string `josn:"client_secret" binding:"required"`
}

type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
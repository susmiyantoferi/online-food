package dto

import "github.com/golang-jwt/jwt/v5"

type TokenClaim struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Username     string `json:"username"`
	Token        string `json:"access_token"`
	TokenRefresh string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExipresIn    int    `json:"expires_in"`
}

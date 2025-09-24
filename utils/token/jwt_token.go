package token

import (
	"errors"
	"fmt"
	"online-food/dto"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId uint, username, email, role string, exp time.Duration) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtExp := time.Now().Add(exp * time.Hour)

	tokenCLaim := &dto.TokenClaim{
		UserID:   userId,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtExp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenCLaim)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

func ClaimTokenRefresh(tokenUser string) (*dto.TokenClaim, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claim := &dto.TokenClaim{}

	token, err := jwt.ParseWithClaims(tokenUser, claim, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")

	}

	return claim, nil
}

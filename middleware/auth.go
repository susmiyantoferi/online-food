package middleware

import (
	"fmt"
	"net/http"
	"online-food/dto"
	"online-food/utils/response"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response.ToResponseJson(ctx, http.StatusUnauthorized, "Unauthorization", "need authorization header", nil)
			ctx.Abort()
			return
		}

		tokenPart := strings.Split(authHeader, " ")
		if len(tokenPart) != 2 || tokenPart[0] != "Bearer" {
			response.ToResponseJson(ctx, http.StatusUnauthorized, "Unauthorization", "invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenStr := tokenPart[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		claim := &dto.TokenClaim{}

		token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.ToResponseJson(ctx, http.StatusUnauthorized, "Unauthorization", "invalid token", nil)
			ctx.Abort()
			return
		}

		ctx.Set("user", claim)
		ctx.Next()
	}

}

func RoleAccessMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userClaims, exist := ctx.Get("user")
		if !exist {
			response.ToResponseJson(ctx, http.StatusUnauthorized, "Unauthorization", "user not found", nil)
			ctx.Abort()
			return
		}

		user := userClaims.(*dto.TokenClaim)
		role := user.Role

		for _, v := range allowedRoles {
			if role == v {
				ctx.Next()
				return
			}
		}

		response.ToResponseJson(ctx, http.StatusForbidden, "Forbidden", "role no permission", nil)
		ctx.Abort()
	}
}

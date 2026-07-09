package jwt

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	golangjwt "github.com/golang-jwt/jwt/v5"

	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

var jwtSecret = []byte("personal-site-secret")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	golangjwt.RegisteredClaims
}

func GenerateToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: golangjwt.RegisteredClaims{
			ExpiresAt: golangjwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  golangjwt.NewNumericDate(time.Now()),
		},
	}

	token := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := golangjwt.ParseWithClaims(tokenString, &Claims{}, func(token *golangjwt.Token) (interface{}, error) {
		_, ok := token.Method.(*golangjwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware() app.HandlerFunc {
	return func(_ context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.WriteErrorMessage(c, errno.Unauthorized, "missing authorization header")
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			response.WriteErrorMessage(c, errno.Unauthorized, "invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)
		claims, err := ParseToken(tokenString)
		if err != nil {
			response.WriteErrorMessage(c, errno.Unauthorized, "invalid token")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next(context.Background())
	}
}

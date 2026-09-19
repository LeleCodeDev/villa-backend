package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/pkg/errors"
)

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.JWTSecret))
}

func ExtractToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("Invalid signind method!")
		}
		return []byte(config.Env.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Unauthorized("Invalid token claims!")
	}

	return claims, nil
}

func ExtractUserID(claims jwt.MapClaims) (uint, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.Unauthorized("Invalid user id in token!")
	}

	return uint(userID), nil
}

package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims расширенная структура из github.com/golang-jwt/jwt/v5 на UserID.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// NewJWTToken создает новый JWT-токен и подписывает его.
func NewJWTToken(
	expiresAt time.Duration,
	secretKey string,
	userID string,
) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			},
			UserID: userID,
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

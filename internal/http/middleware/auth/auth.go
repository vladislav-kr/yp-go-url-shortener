package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	jwttoken "github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware/auth/jwt-token"
)

const key = "auth-token"

// Auth авторизация пользователей.
type Auth struct {
	secretKey string
}

// New конструктор Auth.
func New(secretKey string) *Auth {
	return &Auth{
		secretKey: secretKey,
	}
}

// CreateCookie создает cookie с подписанным токеном JWT.
func (a *Auth) CreateCookie(
	expiresAt time.Duration,
	userID string,
) (*http.Cookie, error) {
	token, err := jwttoken.NewJWTToken(expiresAt, a.secretKey, userID)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:   key,
		Value:  token,
		Path:   "/",
		MaxAge: int(expiresAt.Seconds()),
	}, nil
}

// Validate проверяет токен в cookie.
func (a *Auth) Validate(tokenString string) (*jwttoken.Claims, error) {
	claims := &jwttoken.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.secretKey), nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	if len(claims.UserID) == 0 {
		return nil, fmt.Errorf("user is not found")
	}
	return claims, nil
}

// CookieFromRequest cookie из *http.Request.
func CookieFromRequest(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(key)
}

type contextKey struct {
	name string
}

var (
	userIDCtxKey = &contextKey{"userID"}
)

// ContextWithUserID контекст с UserID.
func ContextWithUserID(parent context.Context, userID string) context.Context {
	ctx := context.WithValue(parent, userIDCtxKey, userID)
	return ctx
}

// UserIDFromContext UserID из контекста.
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDCtxKey).(string); ok {
		return userID
	}
	return ""
}

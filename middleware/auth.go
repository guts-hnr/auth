package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	jwtSecret string
}

func NewAythMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		userID, err := utils.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			userIDKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

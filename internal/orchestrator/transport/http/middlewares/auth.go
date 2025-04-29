package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/bulbosaur/calculator-with-authorization/config"
	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
)

// AuthMiddleware обеспечивает JWT-аутентификацию для API
func AuthMiddleware(cfg *config.JWTConfig, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseJWT(tokenParts[1], cfg.SecretKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

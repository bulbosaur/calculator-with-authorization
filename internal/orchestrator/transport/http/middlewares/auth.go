package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/gorilla/mux"
)

// AuthMiddleware обеспечивает JWT-аутентификацию для API
func AuthMiddleware(authService *auth.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				token := r.URL.Query().Get("token")
				if token == "" {
					http.Error(w, "Missing token", http.StatusUnauthorized)
					return
				}
				authHeader = "Bearer " + token
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			claims, err := authService.ParseJWT(tokenParts[1])
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), models.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

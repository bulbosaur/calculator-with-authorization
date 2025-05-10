package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/mock"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 123}, nil
		},
	}

	middleware := middlewares.AuthMiddleware(mockService)
	testHandler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(models.UserIDKey)
		assert.Equal(t, 123, userID)
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/api/v1/secure", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return nil, fmt.Errorf("invalid token")
		},
	}

	middleware := middlewares.AuthMiddleware(mockService)
	testHandler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/api/v1/secure", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockService := &mock.AuthProvider{}

	middleware := middlewares.AuthMiddleware(mockService)
	testHandler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/api/v1/secure", nil)

	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing token")
}

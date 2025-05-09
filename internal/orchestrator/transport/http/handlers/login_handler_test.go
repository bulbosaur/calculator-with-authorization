package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/mock"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func TestLoginHandler_InvalidRequestBody(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	Service := &auth.Service{
		SecretKey:     "testsecret",
		TokenDuration: time.Hour,
	}
	handler := handlers.LoginHandler(Service, exprRepo)

	req, _ := http.NewRequest("POST", "/api/v1/login", io.NopCloser(bytes.NewBufferString("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, w.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Error != "Bad request" || response.ErrorMessage != models.ErrorInvalidRequestBody.Error() {
		t.Errorf("Unexpected error response: %v", response)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	Service := &auth.Service{
		SecretKey:     "testsecret",
		TokenDuration: time.Hour,
	}
	handler := handlers.LoginHandler(Service, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("nonexistent_user").
		WillReturnError(sql.ErrNoRows)

	reqBody := `{"login":"nonexistent_user","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/v1/login", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(reqBody))

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, w.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Error != "Unauthorized" || response.ErrorMessage != "user not found" {
		t.Errorf("Unexpected error response: %v", response)
	}
}

func TestLoginHandler_InvalidPassword(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	Service := &auth.Service{
		SecretKey:     "testsecret",
		TokenDuration: time.Hour,
	}
	handler := handlers.LoginHandler(Service, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("existing_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).
			AddRow(1, "existing_user", "$2a$10$examplehashedpassword"))

	reqBody := `{"login":"existing_user","password":"wrong_password"}`
	req, _ := http.NewRequest("POST", "/api/v1/login", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(reqBody))

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, w.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Error != "Unauthorized" || response.ErrorMessage != "Invalid password" {
		t.Errorf("Unexpected error response: %v", response)
	}
}

func TestLoginHandler_SuccessfulLogin(t *testing.T) {
	mockAuth := &mock.AuthProvider{
		CompareFunc: func(hash, password string) bool {
			return true
		},
		GenerateJWTFunc: func(userID int) (string, error) {
			claims := &auth.Claims{
				UserID: userID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			return token.SignedString([]byte("testsecret"))
		},
	}

	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.LoginHandler(mockAuth, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("existing_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).
			AddRow(1, "existing_user", "$2a$10$examplehashedpassword"))

	reqBody := `{"login":"existing_user","password":"correct_password"}`
	req := httptest.NewRequest("POST", "/api/v1/login", io.NopCloser(bytes.NewBufferString(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	viper.Set("jwt.secret_key", "testsecret")
	viper.Set("jwt.token_duration", 1)

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if token, ok := response["token"]; !ok || token == "" {
		t.Errorf("Token is missing in response")
	}

	if message, ok := response["message"]; !ok || message != "Authentication successful" {
		t.Errorf("Unexpected message: %v", message)
	}
}

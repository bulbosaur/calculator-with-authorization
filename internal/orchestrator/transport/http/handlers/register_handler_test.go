package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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
)

type MockPasswordHasher struct{}

func (m *MockPasswordHasher) GenerateHash(password string) (string, error) {
	return "", errors.New("mock error")
}

func (m *MockPasswordHasher) Compare(hash, password string) bool {
	return hash == "valid_hash" && password == "valid_password"
}

func TestRegister_InvalidRequestBody(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	authService := &auth.AuthService{
		SecretKey:     "testsecret",
		TokenDuration: time.Hour,
	}
	handler := handlers.Register(authService, exprRepo)

	req, _ := http.NewRequest("POST", "/api/v1/register", io.NopCloser(bytes.NewBufferString("")))
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

func TestRegister_UserAlreadyExists(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	authService := &auth.AuthService{
		SecretKey:     "testsecret",
		TokenDuration: time.Hour,
	}
	handler := handlers.Register(authService, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("existing_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password_hash"}).
			AddRow(1, "existing_user", "hash"))

	reqBody := `{"login":"existing_user","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/v1/register", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(reqBody))

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d; got %d", http.StatusConflict, w.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Error != "Conflict" || response.ErrorMessage != "User already exists" {
		t.Errorf("Unexpected error response: %v", response)
	}
}

func TestRegister_PasswordHashGenerationError(t *testing.T) {
	mockAuth := &mock.MockAuthProvider{
		GenerateHashFunc: func(password string) (string, error) {
			return "", errors.New("mock hash generation error")
		},
	}
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.Register(mockAuth, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("new_user").
		WillReturnError(sql.ErrNoRows)

	reqBody := `{"login":"new_user","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/register", io.NopCloser(bytes.NewBufferString(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d; got %d", http.StatusInternalServerError, w.Code)
	}
	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if response.ErrorMessage != "Failed to generate password hash" {
		t.Errorf("Unexpected error message: %s", response.ErrorMessage)
	}
}

func TestRegister_CreateUserError(t *testing.T) {
	mockAuth := &mock.MockAuthProvider{
		GenerateHashFunc: func(password string) (string, error) {
			return "hashed_password", nil
		},
	}

	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.Register(mockAuth, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("new_user").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users \\(login, password_hash\\) VALUES \\(\\?, \\?\\)").
		WithArgs("new_user", "hashed_password").
		WillReturnError(errors.New("database error"))

	reqBody := `{"login":"new_user","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/v1/register", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(reqBody))

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d; got %d", http.StatusInternalServerError, w.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Error != "Internal error" || response.ErrorMessage != "Failed to create user" {
		t.Errorf("Unexpected error response: %v", response)
	}
}

func TestRegister_SuccessfulRegistration(t *testing.T) {
	mockAuth := &mock.MockAuthProvider{
		GenerateHashFunc: func(password string) (string, error) {
			return "hashed_password", nil
		},
	}
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.Register(mockAuth, exprRepo)

	mock.ExpectQuery("SELECT id, login, password_hash FROM users WHERE login = \\?").
		WithArgs("new_user").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users \\(login, password_hash\\) VALUES \\(\\?, \\?\\)").
		WithArgs("new_user", "hashed_password").
		WillReturnResult(sqlmock.NewResult(1, 1))

	reqBody := `{"login":"new_user","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/register", io.NopCloser(bytes.NewBufferString(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, w.Code)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if message, ok := response["message"].(string); !ok || message != "User created successfully" {
		t.Errorf("Unexpected message: %v", message)
	}
	if userID, ok := response["user_id"].(float64); !ok || int(userID) != 1 {
		t.Errorf("Unexpected user ID: %v", userID)
	}
}

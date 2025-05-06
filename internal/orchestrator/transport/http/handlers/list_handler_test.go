package handlers_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func setup() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return db, mock, err
}

func TestListHandler_MissingAuthHeader(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	req, _ := http.NewRequest("GET", "/expressions", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListHandler_InvalidAuthFormat(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	req, _ := http.NewRequest("GET", "/expressions", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListHandler_InvalidToken(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	req, _ := http.NewRequest("GET", "/expressions", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	viper.Set("jwt.secret_key", "testsecret")

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListHandler_DBError(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	claims := &auth.Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("testsecret"))

	req, _ := http.NewRequest("GET", "/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT \\* FROM expressions WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	viper.Set("jwt.secret_key", "testsecret")

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d; got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestListHandler_Success(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	claims := &auth.Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("testsecret"))

	req, _ := http.NewRequest("GET", "/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()

	rows := sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "error_message"}).
		AddRow(1, 1, "2+2", "completed", "4", "").
		AddRow(2, 1, "5/0", "failed", "", "division by zero")

	mock.ExpectQuery("SELECT \\* FROM expressions WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	viper.Set("jwt.secret_key", "testsecret")

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	var response []models.Expression
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 expressions; got %d", len(response))
	}

	if response[0].Expression != "2+2" || response[1].Expression != "5/0" {
		t.Errorf("Unexpected expressions in response")
	}
}

func TestListHandler_EmptyResult(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.ListHandler(exprRepo)

	claims := &auth.Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("testsecret"))

	req, _ := http.NewRequest("GET", "/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()

	rows := sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "error_message"})

	mock.ExpectQuery("SELECT \\* FROM expressions WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	viper.Set("jwt.secret_key", "testsecret")

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body == "null" {
		t.Errorf("Expected empty array response; got %s", body)
	}
}

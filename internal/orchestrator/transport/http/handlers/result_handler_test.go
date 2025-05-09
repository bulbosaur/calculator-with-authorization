package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/mock"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestResultHandler_MissingAuthHeader(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	mockService := &mock.AuthProvider{}

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization header")
}

func TestResultHandler_InvalidAuthFormat(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	mockService := &mock.AuthProvider{}

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization format")
}

func TestResultHandler_InvalidToken(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return nil, jwt.ErrInvalidKey
		},
	}

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestResultHandler_InvalidExpressionID(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1}, nil
		},
	}

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/invalid_id", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response.ErrorMessage, "strconv.Atoi")
}

func TestResultHandler_RepoError(t *testing.T) {
	db, mockDB, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}

	mockDB.ExpectQuery("SELECT id, user_id, expression, status, result, error_message FROM expressions WHERE id = \\?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1}, nil
		},
	}

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&response)

	assert.Contains(t, response.Error, "ask receiving error")
}

func TestResultHandler_Forbidden(t *testing.T) {
	db, mockDB, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}

	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1}, nil
		},
	}

	expression := &models.Expression{
		ID:           1,
		UserID:       2,
		Expression:   "2+2",
		Status:       "completed",
		Result:       4,
		ErrorMessage: "",
	}

	mockDB.ExpectQuery("SELECT id, user_id, expression, status, result, error_message FROM expressions WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "error_message"}).
			AddRow(expression.ID, expression.UserID, expression.Expression, expression.Status, expression.Result, expression.ErrorMessage))

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestResultHandler_Success(t *testing.T) {
	db, mockDB, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}

	mockService := &mock.AuthProvider{
		ParseJWTFunc: func(tokenString string) (*auth.Claims, error) {
			return &auth.Claims{UserID: 1}, nil
		},
	}

	expression := &models.Expression{
		ID:           1,
		UserID:       1,
		Expression:   "2+2",
		Status:       "completed",
		Result:       4,
		ErrorMessage: "",
	}

	mockDB.ExpectQuery("SELECT id, user_id, expression, status, result, error_message FROM expressions WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "error_message"}).
			AddRow(expression.ID, expression.UserID, expression.Expression, expression.Status, expression.Result, expression.ErrorMessage))

	handler := handlers.ResultHandler(mockService, exprRepo)

	req, _ := http.NewRequest("GET", "/api/v1/result/1", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Response
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, *expression, response.Expression)
}

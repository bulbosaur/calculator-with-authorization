package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRegHandler_InvalidRequestBody(t *testing.T) {
	db, _, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}
	handler := handlers.RegHandler(exprRepo)
	req := httptest.NewRequest("POST", "/register", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Bad request")
}

func TestRegHandler_DBError(t *testing.T) {
	db, mock, _ := setup()
	exprRepo := &repository.ExpressionModel{DB: db}

	mock.ExpectExec("INSERT INTO expressions").
		WithArgs(sqlmock.AnyArg(), "2+2", models.StatusWait, 0.0).
		WillReturnError(errors.New("DB error"))

	handler := handlers.RegHandler(exprRepo)

	reqBody := `{"expression":"2+2"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	userID := 1
	ctx := context.WithValue(req.Context(), models.UserIDKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "something went wrong")
}

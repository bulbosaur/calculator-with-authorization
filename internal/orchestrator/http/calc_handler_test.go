package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

func TestRegHandler(t *testing.T) {
	db, err := repository.InitDB(viper.GetString("database.DATABASE_PATH"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	exprRepo := repository.NewExpressionModel(db)

	cases := []struct {
		requestBody   models.Request
		wantStatus    int
		wantResponse  models.RegisteredExpression
		wantError     models.ErrorResponse
		wantDBRecords int
	}{
		{
			requestBody: models.Request{
				Expression: "2+2",
			},
			wantStatus:    http.StatusCreated,
			wantResponse:  models.RegisteredExpression{ID: 1},
			wantError:     models.ErrorResponse{},
			wantDBRecords: 1,
		},
		{
			requestBody: models.Request{
				Expression: "2++2",
			},
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: models.RegisteredExpression{},
			wantError: models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: models.ErrorInvalidInput.Error(),
			},
			wantDBRecords: 1,
		},
		{
			requestBody: models.Request{
				Expression: "1/0",
			},
			wantStatus:    http.StatusCreated,
			wantResponse:  models.RegisteredExpression{ID: 1},
			wantError:     models.ErrorResponse{},
			wantDBRecords: 2,
		},
		{
			requestBody: models.Request{
				Expression: "2+a",
			},
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: models.RegisteredExpression{},
			wantError: models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: models.ErrorInvalidCharacter.Error(),
			},
			wantDBRecords: 2,
		},
		{
			requestBody: models.Request{
				Expression: "",
			},
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: models.RegisteredExpression{},
			wantError: models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: models.ErrorEmptyExpression.Error(),
			},
			wantDBRecords: 2,
		},
	}

	for _, tc := range cases {
		t.Run("TestRegHandler", func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(regHandler(exprRepo))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("regHandler returned %v, but want %v", status, tc.wantStatus)
			}

			if tc.wantStatus == http.StatusOK {
				var response models.RegisteredExpression
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if tc.wantDBRecords == 1 {
					var expr models.Expression
					err := db.QueryRow("SELECT id, expression, status FROM expressions WHERE id = ?", response.ID).Scan(&expr.ID, &expr.Expression, &expr.Status)
					if err != nil {
						t.Fatalf("Failed to check db record for id: %v, err: %v", response.ID, err)
					}

					if expr.ID != response.ID {
						t.Fatalf("expression id does not match with database. %v != %v", response.ID, expr.ID)
					}
				}

			} else if tc.wantStatus == http.StatusUnprocessableEntity {
				var errorResponse models.ErrorResponse
				err = json.NewDecoder(rr.Body).Decode(&errorResponse)
				if err != nil {
					t.Fatalf("failed to decode error response body: %v", err)
				}

				if errorResponse.Error != tc.wantError.Error {
					t.Errorf("regHandler returned error %v, but want %v", errorResponse.Error, tc.wantError.Error)
				}

				if errorResponse.ErrorMessage != tc.wantError.ErrorMessage {
					t.Errorf("regHandler returned error message %v, but want %v", errorResponse.ErrorMessage, tc.wantError.ErrorMessage)
				}

				if tc.wantDBRecords == 1 {
					var expr models.Expression
					err := db.QueryRow("SELECT id, expression, status FROM expressions WHERE expression = ?", tc.requestBody.Expression).Scan(&expr.ID, &expr.Expression, &expr.Status)
					if err != nil {
						t.Fatalf("Failed to check db record for expression: %v, err: %v", tc.requestBody.Expression, err)
					}
					if expr.Expression != tc.requestBody.Expression {
						t.Fatalf("expression does not match with database. %v != %v", tc.requestBody.Expression, expr.Expression)
					}
				}
			}

			_, err = db.Exec("DELETE FROM expressions")
			if err != nil {
				t.Fatalf("failed to delete db records. %v", err)
			}

			_, err = db.Exec("DELETE FROM tasks")
			if err != nil {
				t.Fatalf("failed to delete db records. %v", err)
			}
		})
	}
}

func TestRegHandlerInvalidJSON(t *testing.T) {
	db, err := repository.InitDB(viper.GetString("database.DATABASE_PATH"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()
	exprRepo := repository.NewExpressionModel(db)

	invalidJSON := []byte(`{"expression": "2+2"`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(regHandler(exprRepo))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned %v, but want %v", status, http.StatusBadRequest)
	}

	var errorResponse models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to decode error response body: %v", err)
	}
	wantError := "Bad request"
	if errorResponse.Error != wantError {
		t.Errorf("Handler returned %v, but want %v", errorResponse.Error, wantError)
	}
	if errorResponse.ErrorMessage != models.ErrorInvalidRequestBody.Error() {
		t.Errorf("Handler returned %v, but want %v", errorResponse.ErrorMessage, models.ErrorInvalidRequestBody.Error())
	}
}

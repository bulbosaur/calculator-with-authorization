package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/spf13/viper"
)

func TestListHandler(t *testing.T) {
	db, err := repository.InitDB(viper.GetString("database.DATABASE_PATH"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	exprRepo := repository.NewExpressionModel(db)
	_, err = db.Exec("INSERT INTO expressions (expression, status, result, error_message) VALUES ('2+2', 'done', '4', '')")
	if err != nil {
		t.Fatalf("failed to insert db records. %v", err)
	}
	_, err = db.Exec("INSERT INTO expressions (expression, status, result, error_message) VALUES ('3*3', 'failed', '', 'Invalid operand')")
	if err != nil {
		t.Fatalf("failed to insert db records. %v", err)
	}
	_, err = db.Exec("INSERT INTO expressions (expression, status, result, error_message) VALUES ('1/0', 'done', '0', 'division by zero')")
	if err != nil {
		t.Fatalf("failed to insert db records. %v", err)
	}

	cases := []struct {
		expectedStatus   int
		expectedResponse []models.Expression
		expectedError    string
	}{
		{
			expectedStatus: http.StatusOK,
			expectedResponse: []models.Expression{
				{ID: 1, Expression: "2+2", Status: "done", Result: 4, ErrorMessage: ""},
				{ID: 2, Expression: "3*3", Status: "failed", Result: 0, ErrorMessage: "Invalid operand"},
				{ID: 3, Expression: "1/0", Status: "done", Result: 0, ErrorMessage: "division by zero"},
			},
			expectedError: "",
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(listHandler(exprRepo))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("returned %v, want %v", status, tc.expectedStatus)
			}

			if tc.expectedStatus == http.StatusOK {
				var actualResponse []models.Expression
				err = json.NewDecoder(rr.Body).Decode(&actualResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				for i, expr := range actualResponse {
					if expr.Expression != tc.expectedResponse[i].Expression || expr.Status != tc.expectedResponse[i].Status || expr.Result != tc.expectedResponse[i].Result || expr.ErrorMessage != tc.expectedResponse[i].ErrorMessage {
						t.Errorf("returned %v, want %v", expr, tc.expectedResponse[i])
					}
				}
			} else {
				if rr.Body.String() != tc.expectedError {
					t.Errorf("returned %v, want %v", rr.Body.String(), tc.expectedError)
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

func TestListHandlerQueryFail(t *testing.T) {
	db, err := repository.InitDB(viper.GetString("database.DATABASE_PATH"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE expressions")
	if err != nil {
		t.Fatalf("failed to delete table. %v", err)
	}

	exprRepo := repository.NewExpressionModel(db)
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listHandler(exprRepo))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("returned %v, want %v", status, http.StatusInternalServerError)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS expressions (id INTEGER PRIMARY KEY AUTOINCREMENT, expression TEXT NOT NULL, status TEXT NOT NULL, result TEXT, error_message TEXT DEFAULT \"\")")
	if err != nil {
		t.Fatalf("failed to create table. %v", err)
	}

	_, err = db.Exec("DELETE FROM expressions")
	if err != nil {
		t.Fatalf("failed to delete db records. %v", err)
	}
	_, err = db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("failed to delete db records. %v", err)
	}
}

func TestListHandlerInvalidFloat(t *testing.T) {
	db, err := repository.InitDB(viper.GetString("database.DATABASE_PATH"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	exprRepo := repository.NewExpressionModel(db)
	_, err = db.Exec("INSERT INTO expressions (expression, status, result, error_message) VALUES ('2+2', 'done', 'qwe', '')")
	if err != nil {
		t.Fatalf("failed to insert db records. %v", err)
	}

	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listHandler(exprRepo))
	handler.ServeHTTP(rr, req)

	var actualResponse []models.Expression
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if actualResponse[0].Result != 0 {
		t.Errorf("returned %v, want %v", actualResponse[0].Result, 0)
	}

	_, err = db.Exec("DELETE FROM expressions")
	if err != nil {
		t.Fatalf("failed to delete db records. %v", err)
	}
	_, err = db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("failed to delete db records. %v", err)
	}
}

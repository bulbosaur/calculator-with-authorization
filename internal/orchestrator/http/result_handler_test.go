package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/gorilla/mux"
)

func TestResultHandlerWithoutMocks(t *testing.T) {
	expressionID := 1
	expression := models.Expression{
		ID:         expressionID,
		Expression: "3+3",
		Result:     6,
	}

	storage := make(map[int]models.Expression)
	storage[expressionID] = expression

	resultHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid expression ID", http.StatusBadRequest)
			return
		}

		expr, exists := storage[id]
		if !exists {
			http.Error(w, "Expression not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]models.Expression{"Expression": expr})
	}

	tests := []struct {
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"Expression\":{\"id\":1,\"expression\":\"3+3\",\"status\":\"\",\"result\":6,\"ErrorMessage\":\"\"}}\n",
		},
		{
			id:             "2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Expression not found\n",
		},
		{
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid expression ID\n",
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", "/expression/"+test.id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/expression/{id}", resultHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
		}

		if rr.Body.String() != test.expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), test.expectedBody)
		}
	}
}

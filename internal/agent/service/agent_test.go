package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

func TestGetTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := models.Task{
			ID:          1,
			Arg1:        10,
			Arg2:        20,
			Operation:   "+",
			PrevTaskID1: 0,
			PrevTaskID2: 0,
		}
		resp := models.TaskResponse{Task: task}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	task, err := getTask(ts.URL)

	if err != nil {
		t.Fatalf("Expected nil error, but got: %v", err)
	}
	if task.ID != 1 {
		t.Errorf("Expexted task.ID = 1, but got: %d", task.ID)
	}
	if task.Arg1 != 10.0 {
		t.Errorf("Expexted task.Arg1 = 10.0, but got: %f", task.Arg1)
	}
	if task.Arg2 != 20.0 {
		t.Errorf("Expexted task.Arg2 = 20.0, but got: %f", task.Arg2)
	}
	if task.Operation != "+" {
		t.Errorf("Expexted task.Operation = '+', but got: %s", task.Operation)
	}
}

func TestExecuteTask(t *testing.T) {
	viper.Set("duration.TIME_ADDITION_MS", 100)
	viper.Set("duration.TIME_SUBTRACTION_MS", 100)
	viper.Set("duration.TIME_MULTIPLICATIONS_MS", 100)
	viper.Set("duration.TIME_DIVISIONS_MS", 100)

	cases := []struct {
		task      models.Task
		expected  float64
		expectErr bool
	}{
		{
			task:      models.Task{ID: 1, Arg1: 10, Arg2: 20, Operation: "+"},
			expected:  30,
			expectErr: false,
		},
		{
			task:      models.Task{ID: 2, Arg1: 20, Arg2: 10, Operation: "-"},
			expected:  10,
			expectErr: false,
		},
		{
			task:      models.Task{ID: 3, Arg1: 10, Arg2: 20, Operation: "*"},
			expected:  200,
			expectErr: false,
		},
		{
			task:      models.Task{ID: 4, Arg1: 20, Arg2: 10, Operation: "/"},
			expected:  2,
			expectErr: false,
		},
		{
			task:      models.Task{ID: 5, Arg1: 20, Arg2: 0, Operation: "/"},
			expected:  0,
			expectErr: false,
		},
		{
			task:      models.Task{ID: 6, Arg1: 10, Arg2: 20, Operation: "hubabuba"},
			expected:  0,
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			result, _, err := executeTask("http://dummy.url", &tc.task)

			if tc.expectErr {
				if err == nil {
					t.Error("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected nil error, but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expexted %f, bur got: %f", tc.expected, result)
				}
			}
		})
	}
}

func TestGetTaskResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := models.Task{
			ID:     1,
			Result: 42,
			Status: models.StatusResolved,
		}
		resp := models.TaskResponse{Task: task}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	result, err := getTaskResult(ts.URL, 1)

	if err != nil {
		t.Fatalf("Expected nil error, but got: %v", err)
	}

	if result != 42.0 {
		t.Errorf("Expexted result = 42.0, but got: %f", result)
	}
}

func TestSendResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		if req["id"] != float64(1) {
			t.Errorf("Ожидалось req['id'] = 1, но получили %v", req["id"])
		}

		if req["result"] != 42.0 {
			t.Errorf("Ожидалось req['result'] = 42.0, но получили %v", req["result"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := sendResult(ts.URL, 1, 42, "")

	if err != nil {
		t.Fatalf("Ожидалось отсутствие ошибки, но получили: %v", err)
	}
}

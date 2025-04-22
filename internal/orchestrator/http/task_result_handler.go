package orchestrator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/gorilla/mux"
)

func taskResultHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		intID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Invalid task ID",
				ErrorMessage: err.Error(),
			})
			return
		}

		task, err := exprRepo.GetTaskByID(intID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "ask receiving error",
				ErrorMessage: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.TaskResponse{
			Task: *task,
		})
	}
}

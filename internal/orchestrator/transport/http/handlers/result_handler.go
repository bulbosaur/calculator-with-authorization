package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/gorilla/mux"
)

// ResultHandler выводит всю информацию по конкретному выражению
func ResultHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		intID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Invalid expression ID",
				ErrorMessage: err.Error(),
			})
			return
		}

		expr, err := exprRepo.GetExpression(intID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "ask receiving error",
				ErrorMessage: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.Response{
			Expression: *expr,
		})
	}
}

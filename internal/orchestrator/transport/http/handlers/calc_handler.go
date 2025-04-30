package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	orchestrator "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/service"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

// RegHandler принимает введенное пользователем выражение и занимается его дальнейшей обработкой. После всех валидаций и подсчетов возвращает результат и код ответа
func RegHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(models.Request)
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Bad request",
				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
			})
			return
		}

		userID, ok := r.Context().Value(models.UserIDKey).(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		id, err := exprRepo.Insert(request.Expression, userID)
		if err != nil {
			log.Printf("something went wrong while creating a record in the database. %v", err)
			return
		}
		log.Printf("Expression ID-%d has been registered", id)

		err = orchestrator.Calc(request.Expression, id, exprRepo)
		if err != nil {
			exprRepo.UpdateStatus(id, models.StatusFailed)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: err.Error(),
			})
			return
		}

		response := models.RegisteredExpression{
			ID: id,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

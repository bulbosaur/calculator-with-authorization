package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	orchestrator "github.com/bulbosaur/web-calculator-golang/internal/orchestrator/service"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func regHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
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

		id, err := exprRepo.Insert(request.Expression)
		if err != nil {
			log.Printf("something went wrong while creating a record in the database. %v", err)
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

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

// ListHandler выводит список всех выражений
func ListHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expressions []models.Expression
		var result string

		rows, err := exprRepo.DB.Query("SELECT * FROM expressions")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var expr models.Expression
			err := rows.Scan(&expr.ID, &expr.UserID, &expr.Expression, &expr.Status, &result, &expr.ErrorMessage)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			expr.Result, err = strconv.ParseFloat(result, 64)
			if err != nil {
				log.Println(err)
			}

			expressions = append(expressions, expr)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expressions)
	}
}

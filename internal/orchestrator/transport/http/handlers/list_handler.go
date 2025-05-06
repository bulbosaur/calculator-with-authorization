package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/spf13/viper"
)

// ListHandler выводит список всех выражений
func ListHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]

		secretKey := viper.GetString("jwt.secret_key")
		claims, err := auth.ParseJWT(token, secretKey)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID := claims.UserID

		rows, err := exprRepo.DB.Query("SELECT * FROM expressions WHERE user_id = $1", userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var expressions []models.Expression
		var result string

		for rows.Next() {
			var expr models.Expression
			err := rows.Scan(&expr.ID, &expr.UserID, &expr.Expression, &expr.Status, &result, &expr.ErrorMessage)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			expr.Result, _ = strconv.ParseFloat(result, 64)
			expressions = append(expressions, expr)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expressions)
	}
}

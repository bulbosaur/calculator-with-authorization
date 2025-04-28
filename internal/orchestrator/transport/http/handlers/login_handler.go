package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/config"
	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

func LoginHandler(exprRepo *repository.ExpressionModel, cfg *config.JWTConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.User
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Bad request",
				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
			})
			return
		}

		user, err := exprRepo.GetUserByLogin(creds.Login)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Unauthorized",
				ErrorMessage: err.Error(),
			})
			return
		}

		token, err := auth.GenerateJWT(user.ID, cfg.SecretKey, cfg.TokenDuration)

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

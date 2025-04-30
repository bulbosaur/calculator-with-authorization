package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/spf13/viper"
)

// LoginHandler - хендлер авторизации. Пользователь отправляет запрос POST /api/v1/login { "login": , "password": }
// В ответ получае 200+OK и JWT токен
func LoginHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
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

		// if !auth.CheckPasswordHash(creds.PasswordHash, user.PasswordHash) {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		SecretKey := viper.GetString("jwt.secret_key")
		token, err := auth.GenerateJWT(user.ID, SecretKey)

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":   token,
			"message": "Authentication successful",
		})
	}
}

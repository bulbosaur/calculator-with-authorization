package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

// LoginHandler - хендлер авторизации. Пользователь отправляет запрос POST /api/v1/login { "login": , "password": }
// В ответ получае 200+OK и JWT токен
func LoginHandler(authProvider auth.AuthProvider, exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Bad request",
				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
			})
			return
		}

		user, err := exprRepo.GetUserByLogin(request.Login)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Unauthorized",
				ErrorMessage: "user not found",
			})
			return
		}

		if !authProvider.Compare(user.PasswordHash, request.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Unauthorized",
				ErrorMessage: "Invalid password",
			})
			return
		}

		token, err := authProvider.GenerateJWT(user.ID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		log.Printf("%s logged in successfully", request.Login)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":   token,
			"message": "Authentication successful",
		})
	}
}

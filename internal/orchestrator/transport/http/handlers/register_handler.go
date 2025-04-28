package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

// Register - регистрация пользователя. Он отправляет запрос POST /api/v1/register { "login": , "password": }
// В ответ получаем 200+OK (в случае успеха) или ошибку
func Register(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Bad request",
				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
			})
			return
		}

		exist, _ := exprRepo.GetUserByLogin(user.Login)
		if exist != nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		hash, err := auth.GenerateHash(user.PasswordHash)
		if err != nil {
			http.Error(w, "Failed to generate hash", http.StatusInternalServerError)
			return
		}

		user.PasswordHash = hash

		err = exprRepo.CreateUser(&user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

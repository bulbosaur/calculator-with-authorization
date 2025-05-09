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
func Register(authProvider auth.Provider, exprRepo *repository.ExpressionModel) http.HandlerFunc {
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

		exist, _ := exprRepo.GetUserByLogin(request.Login)
		if exist != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Conflict",
				ErrorMessage: "User already exists",
			})
			return
		}

		hash, err := authProvider.GenerateHash(request.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Internal error",
				ErrorMessage: "Failed to generate password hash",
			})
			return
		}

		user := &models.User{
			Login:        request.Login,
			PasswordHash: hash,
		}

		userID, err := exprRepo.CreateUser(user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Internal error",
				ErrorMessage: "Failed to create user",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User created successfully",
			"user_id": userID,
		})
	}
}

package repository

import (
	"database/sql"
	"errors"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// CreateUser вносит в БД данного юзера
func (e *ExpressionModel) CreateUser(user *models.User) error {
	_, err := e.DB.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", user.Login, user.PasswordHash)
	return err
}

// GetUserByLogin возвращает данные юзера по логину
func (e *ExpressionModel) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	query := `SELECT id, login, password_hash FROM users WHERE login = ?`
	err := e.DB.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrorUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

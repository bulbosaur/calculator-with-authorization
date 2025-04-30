package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// CreateUser вносит в БД данного юзера
func (e *ExpressionModel) CreateUser(user *models.User) (int, error) {
	result, err := e.DB.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", user.Login, user.PasswordHash)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorCreatingDatabaseRecord, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorReceivingID, err)
	}
	user.ID = int(id)
	return int(id), nil
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

package repository

import (
	"database/sql"
	"errors"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// UserModel обертывает пул подключения sql.DB
type UserModel struct {
	DB *sql.DB
}

// NewUserModel создает экземпляр UserModel
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// CreateUser вносит в БД данного юзера
func (u *UserModel) CreateUser(user *models.User) error {
	_, err := u.DB.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", user.Login, user.PasswordHash)
	return err
}

// GetUserByLogin возвращает данные юзера по логину
func (u *UserModel) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	query := `SELECT id, login, password_hash FROM users WHERE login = ?`
	err := u.DB.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrorUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

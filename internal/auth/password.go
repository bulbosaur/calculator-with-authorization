package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHash создает bcrypt хэш
func GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	saltedBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

// Compare сравнивает хэш с паролем
func Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

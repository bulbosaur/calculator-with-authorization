package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Service реализует AuthProvider
type Service struct {
	SecretKey     string
	TokenDuration time.Duration
}

// NewService создает экземпляр Service
func NewService(secretKey string, tokenDuration time.Duration) *Service {
	return &Service{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration * time.Hour,
	}
}

// GenerateHash генирирует хэш из данного пароля
func (s *Service) GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare сравнивает хэш с паролем
func (s *Service) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService реализует AuthProvider
type AuthService struct {
	SecretKey     string
	TokenDuration time.Duration
}

// NewAuthService создает экземпляр AuthService
func NewAuthService(secretKey string, tokenDuration time.Duration) *AuthService {
	return &AuthService{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}

// GenerateHash генирирует хэш из данного пароля
func (s *AuthService) GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare сравнивает хэш с паролем
func (s *AuthService) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

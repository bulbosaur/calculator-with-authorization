package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims - данные JWT
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Provider абстрагирует методы хэширования и работы с JWT
type Provider interface {
	GenerateHash(password string) (string, error)
	Compare(hash, password string) bool
	GenerateJWT(userID int) (string, error)
	ParseJWT(tokenString string) (*Claims, error)
}

// GenerateJWT создает новый токен, подписанный секретным ключом
func (s *Service) GenerateJWT(userID int) (string, error) {
	if s.SecretKey == "" {
		return "", errors.New("secret key is empty")
	}
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.TokenDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

// ParseJWT парсит токен в Claims и проверяет, что он вадилен
func (s *Service) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

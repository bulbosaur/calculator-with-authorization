package mock

import "github.com/bulbosaur/calculator-with-authorization/internal/auth"

// AuthProvider - это мок реализации auth.Provider для тестирования
type AuthProvider struct {
	GenerateHashFunc func(password string) (string, error)
	CompareFunc      func(hash, password string) bool
	GenerateJWTFunc  func(userID int) (string, error)
	ParseJWTFunc     func(tokenString string) (*auth.Claims, error)
}

// GenerateHash имитирует метод GenerateHash
func (m *AuthProvider) GenerateHash(password string) (string, error) {
	if m.GenerateHashFunc != nil {
		return m.GenerateHashFunc(password)
	}
	return "", nil
}

// Compare имитирует метод Compare
func (m *AuthProvider) Compare(hash, password string) bool {
	if m.CompareFunc != nil {
		return m.CompareFunc(hash, password)
	}
	return true
}

// GenerateJWT имитирует метод GenerateJWT
func (m *AuthProvider) GenerateJWT(userID int) (string, error) {
	if m.GenerateJWTFunc != nil {
		return m.GenerateJWTFunc(userID)
	}
	return "mock_token", nil
}

// ParseJWT имитирует метод ParseJWT
func (m *AuthProvider) ParseJWT(tokenString string) (*auth.Claims, error) {
	if m.ParseJWTFunc != nil {
		return m.ParseJWTFunc(tokenString)
	}
	return &auth.Claims{UserID: 1}, nil
}

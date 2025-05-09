package mock

import "github.com/bulbosaur/calculator-with-authorization/internal/auth"

type MockAuthProvider struct {
	GenerateHashFunc func(password string) (string, error)
	CompareFunc      func(hash, password string) bool
	GenerateJWTFunc  func(userID int) (string, error)
	ParseJWTFunc     func(tokenString string) (*auth.Claims, error)
}

func (m *MockAuthProvider) GenerateHash(password string) (string, error) {
	if m.GenerateHashFunc != nil {
		return m.GenerateHashFunc(password)
	}
	return "", nil
}

func (m *MockAuthProvider) Compare(hash, password string) bool {
	if m.CompareFunc != nil {
		return m.CompareFunc(hash, password)
	}
	return true
}

func (m *MockAuthProvider) GenerateJWT(userID int) (string, error) {
	if m.GenerateJWTFunc != nil {
		return m.GenerateJWTFunc(userID)
	}
	return "mock_token", nil
}

func (m *MockAuthProvider) ParseJWT(tokenString string) (*auth.Claims, error) {
	if m.ParseJWTFunc != nil {
		return m.ParseJWTFunc(tokenString)
	}
	return &auth.Claims{UserID: 1}, nil
}

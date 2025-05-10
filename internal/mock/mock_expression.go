package mock

import (
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// MockExpressionModel — минимальная реализация ExpressionModel для тестов
type MockExpressionModel struct {
	models.ExpressionRepository
}

// GetTask — заглушка, возвращающая nil и ошибку
func (m *MockExpressionModel) GetTask() (*models.Task, int, error) {
	return nil, 0, nil
}

// UpdateTaskResult — заглушка
func (m *MockExpressionModel) UpdateTaskResult(id int, result float64, err string) error {
	return nil
}

// GetExpression — заглушка
func (m *MockExpressionModel) GetExpression(id int) (*models.Expression, error) {
	return &models.Expression{}, nil
}

// Insert — заглушка
func (m *MockExpressionModel) Insert(expr string, userID int) (int, error) {
	return 1, nil
}

// UpdateStatus — заглушка
func (m *MockExpressionModel) UpdateStatus(id int, status string) {
}

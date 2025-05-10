package mock

import (
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// ExpressionModel — минимальная реализация ExpressionModel для тестов
type ExpressionModel struct {
	models.ExpressionRepository
}

// GetTask — заглушка, возвращающая nil и ошибку
func (m *ExpressionModel) GetTask() (*models.Task, int, error) {
	return nil, 0, nil
}

// UpdateTaskResult — заглушка
func (m *ExpressionModel) UpdateTaskResult(id int, result float64, err string) error {
	return nil
}

// GetExpression — заглушка
func (m *ExpressionModel) GetExpression(id int) (*models.Expression, error) {
	return &models.Expression{}, nil
}

// Insert — заглушка
func (m *ExpressionModel) Insert(expr string, userID int) (int, error) {
	return 1, nil
}

// UpdateStatus — заглушка
func (m *ExpressionModel) UpdateStatus(id int, status string) {
}

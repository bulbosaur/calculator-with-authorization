package orchestrator

import (
	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

// Calc вызывает токенизацию выражения, записывает его в RPN. а затем в параллельных горутинах подсчитывает значения выражений в скобках
func Calc(stringExpression string, id int, taskRepo *repository.ExpressionModel) error {
	taskRepo.Mu.Lock()
	defer taskRepo.Mu.Unlock()

	expression, err := tokenize(stringExpression)
	if err != nil {
		return err
	}

	if len(expression) == 0 {
		return models.ErrorEmptyExpression
	}

	reversePolishNotation, err := toReversePolishNotation(expression)
	if err != nil {
		return err
	}

	parseRPN(reversePolishNotation, id, taskRepo)

	return nil
}

// NewTask создает экземпляр структуры Task
func NewTask(id int, arg1, arg2 float64, operation string) *models.Task {
	newTask := models.Task{
		ExpressionID: id,
		Arg1:         arg1,
		Arg2:         arg2,
		Operation:    operation,
		Status:       models.StatusNew,
		Result:       0,
	}
	return &newTask
}

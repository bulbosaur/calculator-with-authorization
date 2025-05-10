package repository_test

import (
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestInsertAndGetTask(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("3 + 4", 1)

	task := &models.Task{
		ExpressionID: exprID,
		Arg1:         3,
		Arg2:         4,
		Operation:    "+",
		Status:       models.StatusWait,
	}

	taskID, err := repo.InsertTask(task)
	assert.NoError(t, err)
	assert.Greater(t, taskID, 0)

	dbTask, err := repo.GetTaskByID(taskID)
	assert.NoError(t, err)
	assert.Equal(t, task.ExpressionID, dbTask.ExpressionID)
	assert.Equal(t, task.Arg1, dbTask.Arg1)
	assert.Equal(t, task.Arg2, dbTask.Arg2)
	assert.Equal(t, task.Operation, dbTask.Operation)
	assert.Equal(t, task.Status, dbTask.Status)
}

func TestUpdateTaskResult(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("6 / 2", 1)

	task := &models.Task{
		ExpressionID: exprID,
		Arg1:         6,
		Arg2:         2,
		Operation:    "/",
		Status:       models.StatusWait,
	}

	taskID, _ := repo.InsertTask(task)

	err := repo.UpdateTaskResult(taskID, 3.0, "")
	assert.NoError(t, err)

	dbTask, _ := repo.GetTaskByID(taskID)
	assert.Equal(t, models.StatusResolved, dbTask.Status)
	assert.Equal(t, 3.0, dbTask.Result)

	expr, _ := repo.GetExpression(exprID)
	assert.Equal(t, models.StatusResolved, expr.Status)
	assert.Equal(t, 3.0, expr.Result)
}

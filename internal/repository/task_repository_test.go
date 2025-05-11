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

func TestUpdateTaskResult_InvalidExpressionID(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	_, err := repo.DB.Exec("INSERT INTO tasks (expressionID, arg1, arg2, operation, status) VALUES (?, ?, ?, ?, ?)",
		999, 1, 2, "+", models.StatusWait)
	assert.NoError(t, err)

	var taskID int
	err = repo.DB.QueryRow("SELECT id FROM tasks WHERE expressionID = ?", 999).Scan(&taskID)
	assert.NoError(t, err)

	err = repo.UpdateTaskResult(taskID, 3, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expression with ID 999 does not exist")
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

func TestGetTask(t *testing.T) {
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
	taskID, _ := repo.InsertTask(task)

	dbTask, dbTaskID, err := repo.GetTask()
	assert.NoError(t, err)
	assert.Equal(t, taskID, dbTaskID)
	assert.Equal(t, task.ExpressionID, dbTask.ExpressionID)
	assert.Equal(t, task.Arg1, dbTask.Arg1)
	assert.Equal(t, task.Arg2, dbTask.Arg2)
	assert.Equal(t, task.Operation, dbTask.Operation)
	assert.Equal(t, models.StatusWait, dbTask.Status)

	dbTask, _ = repo.GetTaskByID(taskID)
	assert.Equal(t, models.StatusInProcess, dbTask.Status)

	dbTask, dbTaskID, err = repo.GetTask()
	assert.NoError(t, err)
	assert.Nil(t, dbTask)
	assert.Zero(t, dbTaskID)
}

func TestGetTaskStatus(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("7 * 8", 1)

	task := &models.Task{
		ExpressionID: exprID,
		Arg1:         7,
		Arg2:         8,
		Operation:    "*",
		Status:       models.StatusResolved,
		Result:       56.0,
	}
	taskID, _ := repo.InsertTask(task)

	status, result, err := repo.GetTaskStatus(taskID)
	assert.NoError(t, err)
	assert.Equal(t, models.StatusResolved, status)
	assert.Equal(t, 56.0, result)
}

func TestUpdateTaskResult_WithError(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("5/0", 1)
	task := &models.Task{
		ExpressionID: exprID,
		Arg1:         5,
		Arg2:         0,
		Operation:    "/",
		Status:       models.StatusWait,
	}
	taskID, _ := repo.InsertTask(task)

	err := repo.UpdateTaskResult(taskID, 0, models.ErrorDivisionByZero.Error())
	assert.NoError(t, err)

	expr, _ := repo.GetExpression(exprID)
	assert.Equal(t, models.StatusFailed, expr.Status)
	assert.Equal(t, models.ErrorDivisionByZero.Error(), expr.ErrorMessage)
}

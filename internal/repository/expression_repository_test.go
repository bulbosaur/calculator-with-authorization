package repository_test

import (
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAreAllTasksCompleted(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("5 + 3", 1)

	task1 := &models.Task{
		ExpressionID: exprID,
		Arg1:         5,
		Arg2:         3,
		Operation:    "+",
		Status:       models.StatusResolved,
	}
	_, _ = repo.InsertTask(task1)

	task2 := &models.Task{
		ExpressionID: exprID,
		Arg1:         2,
		Arg2:         2,
		Operation:    "+",
		Status:       models.StatusWait,
	}
	taskID2, _ := repo.InsertTask(task2)

	completed, err := repo.AreAllTasksCompleted(exprID)
	assert.NoError(t, err)
	assert.False(t, completed)

	repo.UpdateTaskStatus(taskID2, models.StatusResolved)

	completed, err = repo.AreAllTasksCompleted(exprID)
	assert.NoError(t, err)
	assert.True(t, completed)
}

func TestSetResult(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("9 / 3", 1)

	repo.SetResult(exprID, 3.0)

	expr, _ := repo.GetExpression(exprID)
	assert.Equal(t, 3.0, expr.Result)
}

func TestGetTaskStatus_InvalidTaskID(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	status, result, err := repo.GetTaskStatus(999)
	assert.Error(t, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())
	assert.Empty(t, status)
	assert.Zero(t, result)
}

func TestUpdateStatus_NonExistentExpression(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	repo.UpdateStatus(999, models.StatusResolved)

	expr, err := repo.GetExpression(999)
	assert.Error(t, err)
	assert.Nil(t, expr)
}

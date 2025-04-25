package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
)

// InsertTask записывает мат выражение в таблицу БД
func (e *ExpressionModel) InsertTask(task *models.Task) (int, error) {
	query := `
        INSERT INTO tasks (expressionID, arg1, arg2, prev_task_id1, prev_task_id2, operation, status, result)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := e.DB.Exec(
		query,
		task.ExpressionID,
		task.Arg1,
		task.Arg2,
		task.PrevTaskID1,
		task.PrevTaskID2,
		task.Operation,
		task.Status,
		task.Result,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get task ID: %v", err)
	}

	return int(id), nil
}

// GetTask забирает из базы таску для агента
func (e *ExpressionModel) GetTask() (*models.Task, int, error) {
	query := `
        SELECT t.id, t.expressionID, 
               COALESCE(t1.result, t.arg1) AS arg1, 
               COALESCE(t2.result, t.arg2) AS arg2, 
               t.prev_task_id1, t.prev_task_id2, 
               t.operation, t.status, t.result
        FROM tasks t
        LEFT JOIN tasks t1 ON t.prev_task_id1 = t1.id
        LEFT JOIN tasks t2 ON t.prev_task_id2 = t2.id
        WHERE t.status = ?
        AND (t.prev_task_id1 = 0 OR t1.status = ?)
        AND (t.prev_task_id2 = 0 OR t2.status = ?)
        LIMIT 1
    `

	var task models.Task
	err := e.DB.QueryRow(query, models.StatusWait, models.StatusResolved, models.StatusResolved).Scan(
		&task.ID,
		&task.ExpressionID,
		&task.Arg1,
		&task.Arg2,
		&task.PrevTaskID1,
		&task.PrevTaskID2,
		&task.Operation,
		&task.Status,
		&task.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to get task: %v", err)
	}

	_, err = e.DB.Exec("UPDATE tasks SET status = ? WHERE id = ?", models.StatusInProcess, task.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to update task status: %v", err)
	}

	return &task, task.ID, nil
}

// GetTaskByID возвращает из базы данных соответствующую таску
func (e *ExpressionModel) GetTaskByID(taskID int) (*models.Task, error) {
	query := `
        SELECT id, expressionID, arg1, arg2, prev_task_id1, prev_task_id2, operation, status, result
        FROM tasks
        WHERE id = ?
    `

	var task models.Task
	err := e.DB.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.ExpressionID,
		&task.Arg1,
		&task.Arg2,
		&task.PrevTaskID1,
		&task.PrevTaskID2,
		&task.Operation,
		&task.Status,
		&task.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

// GetTaskStatus возвращает статус и ответ таски
func (e *ExpressionModel) GetTaskStatus(taskID int) (string, float64, error) {
	var status string
	var result float64

	err := e.DB.QueryRow("SELECT status, result FROM tasks WHERE id = ?", taskID).Scan(&status, &result)
	if err != nil {
		return "", 0, err
	}
	log.Printf("status %v", status)
	return status, result, nil
}

// UpdateTaskStatus обновляет статус таски
func (e *ExpressionModel) UpdateTaskStatus(taskID int, status string) {
	query := "UPDATE tasks SET status = ? WHERE id = ?"

	_, err := e.DB.Exec(query, status, taskID)
	if err != nil {
		log.Println(err)
	}
}

// UpdateTaskResult обновляет результат таски в базе и если все остальные действия выраженрия выполены, пишет окончательный ответ
func (e *ExpressionModel) UpdateTaskResult(taskID int, result float64, errorMessage string) error {
	_, err := e.DB.Exec(
		"UPDATE tasks SET status = ?, result = ?, error_message = ? WHERE id = ?",
		models.StatusResolved,
		result,
		errorMessage,
		taskID,
	)
	if err != nil {
		return err
	}

	if errorMessage != "" {
		log.Printf("update task ID-%d: %v\nerror message: %v", taskID, result, errorMessage)
	} else {
		log.Printf("update result for task ID-%d: %v", taskID, result)
	}

	var exprID int
	err = e.DB.QueryRow("SELECT expressionID FROM tasks WHERE id = ?", taskID).Scan(&exprID)
	if err != nil {
		return fmt.Errorf("failed to get expression ID: %v", err)
	}

	completed, err := e.AreAllTasksCompleted(exprID)
	if err != nil {
		return fmt.Errorf("failed to check tasks completion: %v", err)
	}

	if completed {
		finalResult, errorMessage, err := e.CalculateExpressionResult(exprID)
		if err != nil {
			return fmt.Errorf("failed to calculate expression result: %v", err)
		}

		err = e.UpdateExpressionResult(exprID, finalResult, errorMessage)
		if err != nil {
			return fmt.Errorf("failed to update expression result: %v", err)
		}
	}

	return nil
}

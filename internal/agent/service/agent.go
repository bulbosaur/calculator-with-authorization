package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

// Workers - переменная, в которой хранится количество одновременно работающих воркеров
var Workers int

// RunAgent запускает агента
func RunAgent() {
	orchost := viper.GetString("server.ORC_HOST")
	orcport := viper.GetString("server.ORC_PORT")
	orchestratorURL := fmt.Sprintf("http://%s:%s", orchost, orcport)

	Workers = viper.GetInt("worker.COMPUTING_POWER")
	if Workers <= 0 {
		Workers = 1
	}

	for i := 1; i <= Workers; i++ {
		log.Printf("Starting worker %d", i)
		go worker(i, orchestratorURL)
	}

	log.Printf("Starting %d workers", Workers)

	select {}
}

func getTask(orchestratorURL string) (*models.Task, error) {
	resp, err := http.Get(orchestratorURL + "/internal/task")
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("orchestrator returned status code %d", resp.StatusCode)
	}

	var res models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}

	if res.Task.ID != 0 {
		log.Printf("Received task: ID=%d, Arg1=%f, Arg2=%f, PrevTaskID1=%d, PrevTaskID2=%d, Operation=%s",
			res.Task.ID, res.Task.Arg1, res.Task.Arg2, res.Task.PrevTaskID1, res.Task.PrevTaskID2, res.Task.Operation)
	}

	return &res.Task, nil
}

func executeTask(orchestratorURL string, task *models.Task) (float64, string, error) {
	if task == nil || task.ID == 0 {
		return 0, "", fmt.Errorf("invalid task: task is nil or has ID 0")
	}

	var arg1, arg2 float64
	var err error

	if task.PrevTaskID1 != 0 {
		arg1, err = getTaskResult(orchestratorURL, task.PrevTaskID1)
		if err != nil {
			return 0, "", fmt.Errorf("failed to get result for PrevTaskID1 (%d): %v", task.PrevTaskID1, err)
		}
	} else {
		arg1 = task.Arg1
	}

	if task.PrevTaskID2 != 0 {
		arg2, err = getTaskResult(orchestratorURL, task.PrevTaskID2)
		if err != nil {
			return 0, "", fmt.Errorf("failed to get result for PrevTaskID2 (%d): %v", task.PrevTaskID2, err)
		}
	} else {
		arg2 = task.Arg2
	}

	if task.Operation == "" {
		return 0, "", fmt.Errorf("invalid operation: operation is empty")
	}

	switch task.Operation {
	case "+":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_ADDITION_MS")) * time.Millisecond)
		return arg1 + arg2, "", nil
	case "-":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_SUBTRACTION_MS")) * time.Millisecond)
		return arg1 - arg2, "", nil
	case "*":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_MULTIPLICATIONS_MS")) * time.Millisecond)
		return arg1 * arg2, "", nil
	case "/":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_DIVISIONS_MS")) * time.Millisecond)
		if arg2 == 0 {
			return 0, models.ErrorDivisionByZero.Error(), nil
		}
		return arg1 / arg2, "", nil
	default:
		return 0, "", fmt.Errorf("invalid operation: %s", task.Operation)
	}
}

func getTaskResult(orchestratorURL string, taskID int) (float64, error) {
	resp, err := http.Get(orchestratorURL + fmt.Sprintf("/internal/task/%d", taskID))
	if err != nil {
		return 0, fmt.Errorf("failed to get task result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("orchestrator returned status code %d", resp.StatusCode)
	}

	var res models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("failed to decode task result: %w", err)
	}

	if res.Task.Status != models.StatusResolved {
		return 0, fmt.Errorf("task with ID %d is not done yet", taskID)
	}
	return res.Task.Result, nil
}

func sendResult(orchestratorURL string, taskID int, result float64, errorMessage string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"id":            taskID,
		"result":        result,
		"error_message": errorMessage,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	resp, err := http.Post(orchestratorURL+"/internal/task", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("payload: %s", payload)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("orchestrator returned status code %d", resp.StatusCode)
	}

	return nil
}

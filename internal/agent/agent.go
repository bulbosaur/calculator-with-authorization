package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Workers - переменная, в которой хранится количество одновременно работающих воркеров
var Workers int

type GRPCAgent struct {
	client proto.TaskServiceClient
	conn   *grpc.ClientConn
}

func newGRPCAgent() (*GRPCAgent, error) {
	orchost := viper.GetString("server.GRPC_HOST")
	orcport := viper.GetString("server.GRPC_PORT")
	orchestratorAddr := fmt.Sprintf("%s:%s", orchost, orcport)

	conn, err := grpc.NewClient(
		orchestratorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024),
			grpc.MaxCallSendMsgSize(10*1024*1024),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to orchestrator: %v", err)
	}

	client := proto.NewTaskServiceClient(conn)

	return &GRPCAgent{
		client: client,
		conn:   conn,
	}, nil
}

// RunAgent запускает агента
func RunAgent() {
	agent, err := newGRPCAgent()
	if err != nil {
		log.Printf("Failed to create agent: %v", err)
		return
	}

	Workers = viper.GetInt("worker.COMPUTING_POWER")
	if Workers <= 0 {
		Workers = 1
	}

	for i := 1; i <= Workers; i++ {
		log.Printf("Starting worker %d", i)
		go agent.worker(i)
	}

	log.Printf("Starting %d workers", Workers)

	select {}
}

func (a *GRPCAgent) getTask(ctx context.Context) (*models.Task, error) {
	resp, err := a.client.ReceiveTask(ctx, &proto.GetTaskRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	task := &models.Task{
		ID:           int(resp.Id),
		ExpressionID: int(resp.ExpressionId),
		Arg1:         resp.Arg1,
		Arg2:         resp.Arg2,
		PrevTaskID1:  int(resp.PrevTask_Id1),
		PrevTaskID2:  int(resp.PrevTask_Id2),
		Operation:    resp.Operation,
	}

	if task.ID != 0 {
		log.Printf("Received task: ID=%d, Arg1=%f, Arg2=%f, PrevTaskID1=%d, PrevTaskID2=%d, Operation=%s",
			task.ID, task.Arg1, task.Arg2, task.PrevTaskID1, task.PrevTaskID2, task.Operation)
	}

	return task, nil
}

func (a *GRPCAgent) executeTask(ctx context.Context, task *models.Task) (float64, string, error) {
	if task == nil || task.ID == 0 {
		return 0, "", fmt.Errorf("invalid task: task is nil or has ID 0")
	}

	var arg1, arg2 float64
	var err error

	if task.PrevTaskID1 != 0 {
		arg1, err = a.fetchTaskResult(ctx, task.PrevTaskID1)
		if err != nil {
			return 0, "", fmt.Errorf("failed to get result for PrevTaskID1 (%d): %v", task.PrevTaskID1, err)
		}
	} else {
		arg1 = task.Arg1
	}

	if task.PrevTaskID2 != 0 {
		arg2, err = a.fetchTaskResult(ctx, task.PrevTaskID2)
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

func (a *GRPCAgent) fetchTaskResult(ctx context.Context, taskID int) (float64, error) {
	resp, err := a.client.GetTaskResult(ctx, &proto.GetTaskResultRequest{TaskId: int32(taskID)})
	if err != nil {
		return 0, fmt.Errorf("failed to get task result: %v", err)
	}

	if !resp.Success {
		return 0, fmt.Errorf("task with ID %d is not done yet", taskID)
	}

	return resp.Result, nil
}

func (a *GRPCAgent) sendResult(ctx context.Context, taskID int, result float64, errorMessage string) error {
	_, err := a.client.SubmitTaskResult(context.Background(), &proto.SubmitTaskResultRequest{
		TaskId:       int32(taskID),
		Result:       result,
		ErrorMessage: errorMessage,
	})

	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}

	return nil
}

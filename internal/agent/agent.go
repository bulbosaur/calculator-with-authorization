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

// GRPCAgent - gRPC-клиент для взаимодействия с оркестратором вычислений
type GRPCAgent struct {
	Client proto.TaskServiceClient
	Conn   *grpc.ClientConn
}

func newGRPCAgent() (*GRPCAgent, error) {
	orchost := viper.GetString("server.GRPC_HOST")
	orcport := viper.GetString("server.GRPC_PORT")

	if orchost == "" {
		orchost = "localhost"
	}
	if orcport == "" {
		orcport = "50051"
	}

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
		Client: client,
		Conn:   conn,
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
		go agent.Worker(i)
	}

	log.Printf("Starting %d workers", Workers)

	select {}
}

func (a *GRPCAgent) getTask(ctx context.Context) (*models.Task, error) {
	resp, err := a.Client.ReceiveTask(ctx, &proto.GetTaskRequest{})
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

	arg1 = task.Arg1
	arg2 = task.Arg2

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

func (a *GRPCAgent) sendResult(ctx context.Context, taskID int, result float64, errorMessage string) error {
	Mu.Lock()
	defer Mu.Unlock()

	_, err := a.Client.SubmitTaskResult(context.Background(), &proto.SubmitTaskResultRequest{
		TaskId:       int32(taskID),
		Result:       result,
		ErrorMessage: errorMessage,
	})

	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}

	return nil
}

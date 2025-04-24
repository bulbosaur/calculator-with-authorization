package orchestrator

import (
	"context"
	"log"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskServer реализует gRPC-сервис для управления задачами
type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	exprRepo *repository.ExpressionModel
}

func newTaskServer(repo *repository.ExpressionModel) *TaskServer {
	return &TaskServer{exprRepo: repo}
}

// ReceiveTask обрабатывает запрос от агента на получение задачи
func (ts *TaskServer) ReceiveTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	// Проверка аутентификации
	// if req.GetCtx() == nil || req.GetCtx().GetAuthToken() == "" {
	//     return nil, status.Error(codes.Unauthenticated, "authentication required")
	// }

	task, id, err := ts.exprRepo.GetTask()
	if err != nil {
		log.Println("Failed to get task:", err)
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	if task == nil {
		return nil, status.Error(codes.NotFound, "no tasks available")
	}

	ts.exprRepo.UpdateTaskStatus(id, models.StatusInProcess)

	return &proto.Task{
		Id:           int32(task.ID),
		ExpressionId: int32(task.ExpressionID),
		Arg1:         task.Arg1,
		Arg2:         task.Arg2,
		PrevTask_Id1: int32(task.PrevTaskID1),
		PrevTask_Id2: int32(task.PrevTaskID2),
		Operation:    task.Operation,
		Status:       task.Status,
		Result:       task.Result,
	}, nil
}

// SubmitTaskResult обрабатывает результат выполнения задачи от агента
func (ts *TaskServer) SubmitTaskResult(ctx context.Context, req *proto.SubmitTaskResultRequest) (*proto.SubmitTaskResultResponse, error) {
	err := ts.exprRepo.UpdateTaskResult(
		int(req.TaskId),
		req.Result,
		req.ErrorMessage,
	)
	if err != nil {
		log.Printf("Failed to update task result: %v", err)
		return &proto.SubmitTaskResultResponse{Success: false},
			status.Errorf(codes.Internal, "failed to update task result: %v", err)
	}

	return &proto.SubmitTaskResultResponse{Success: true}, nil
}

// GetTaskResult возвращает результат выполнения задачи
func (ts *TaskServer) GetTaskResult(ctx context.Context, req *proto.GetTaskResultRequest) (*proto.GetTaskResultResponse, error) {
	task, err := ts.exprRepo.GetTaskByID(int(req.TaskId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	if task.Status != models.StatusResolved {
		return &proto.GetTaskResultResponse{
			Success: false,
		}, nil
	}

	return &proto.GetTaskResultResponse{
		Success: true,
		Result:  task.Result,
	}, nil
}

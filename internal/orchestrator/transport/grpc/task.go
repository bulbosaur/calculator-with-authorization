package orchestrator

import (
	"context"
	"log"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskServer реализует gRPC-сервис для управления задачами
type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	ExprRepo models.ExpressionRepository
}

func newTaskServer(repo models.ExpressionRepository) *TaskServer {
	return &TaskServer{ExprRepo: repo}
}

// ReceiveTask обрабатывает запрос от агента на получение задачи
func (ts *TaskServer) ReceiveTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	task, id, err := ts.ExprRepo.GetTask()
	if err != nil {
		log.Println("Failed to get task:", err)
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	if task == nil {
		return nil, status.Error(codes.NotFound, "no tasks available")
	}

	ts.ExprRepo.UpdateTaskStatus(id, models.StatusInProcess)

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
	if req.TaskId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid task ID")
	}

	err := ts.ExprRepo.UpdateTaskResult(
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

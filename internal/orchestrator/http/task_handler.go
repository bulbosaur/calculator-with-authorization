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

type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	exprRepo *repository.ExpressionModel
}

func NewTaskServer(repo *repository.ExpressionModel) *TaskServer {
	return &TaskServer{exprRepo: repo}
}

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

	ts.exprRepo.UpdateTaskStatus(id, models.StatusCalculate)

	return &proto.Task{
		Id:           int32(task.ID),
		ExpressionID: int32(task.ExpressionID),
		Arg1:         task.Arg1,
		Arg2:         task.Arg2,
		PrevTaskID1:  int32(task.PrevTaskID1),
		PrevTaskID2:  int32(task.PrevTaskID2),
		Operation:    task.Operation,
	}, nil
}

func (s *TaskServer) SubmitTaskResult(ctx context.Context, req *proto.SubmitTaskResultRequest) (*proto.SubmitTaskResultResponse, error) {
	err := s.exprRepo.UpdateTaskResult(
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

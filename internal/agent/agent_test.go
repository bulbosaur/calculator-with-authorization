package agent

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testOrchestratorServer struct {
	proto.UnimplementedTaskServiceServer
	tasks    map[int32]*proto.Task
	results  map[int32]float64
	received chan *proto.SubmitTaskResultRequest
}

func (s *testOrchestratorServer) ReceiveTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	for _, task := range s.tasks {
		return task, nil
	}
	return &proto.Task{}, nil
}

func (s *testOrchestratorServer) SubmitTaskResult(ctx context.Context, req *proto.SubmitTaskResultRequest) (*proto.SubmitTaskResultResponse, error) {
	s.received <- req
	return &proto.SubmitTaskResultResponse{}, nil
}

var lis *bufconn.Listener

func startTestServer() (*grpc.Server, *testOrchestratorServer) {
	lis = bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	testServer := &testOrchestratorServer{
		tasks:    make(map[int32]*proto.Task),
		results:  make(map[int32]float64),
		received: make(chan *proto.SubmitTaskResultRequest, 10),
	}
	proto.RegisterTaskServiceServer(srv, testServer)
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv, testServer
}

func TestExecuteTask(t *testing.T) {
	agent := &GRPCAgent{}

	tests := []struct {
		name       string
		task       *models.Task
		wantResult float64
		wantErrMsg string
		wantErr    bool
	}{
		{
			name:       "Addition",
			task:       &models.Task{ID: 1, Operation: "+", Arg1: 5, Arg2: 3},
			wantResult: 8,
		},
		{
			name:       "DivisionByZero",
			task:       &models.Task{ID: 2, Operation: "/", Arg1: 5, Arg2: 0},
			wantErrMsg: models.ErrorDivisionByZero.Error(),
		},
		{
			name:    "InvalidOperation",
			task:    &models.Task{Operation: "invalid"},
			wantErr: true,
		},
		{
			name:       "Multiplication",
			task:       &models.Task{ID: 3, Operation: "*", Arg1: 4, Arg2: 5},
			wantResult: 20,
		},
		{
			name:       "Subtraction",
			task:       &models.Task{ID: 4, Operation: "-", Arg1: 10, Arg2: 7},
			wantResult: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errMsg, err := agent.executeTask(context.Background(), tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantResult, result)
			assert.Contains(t, errMsg, tt.wantErrMsg)
		})
	}
}

func TestGRPCAgentIntegration(t *testing.T) {
	srv, testServer := startTestServer()
	defer srv.Stop()

	testServer.tasks[1] = &proto.Task{
		Id:        1,
		Arg1:      10,
		Arg2:      20,
		Operation: "+",
	}

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()

	agent := &GRPCAgent{
		client: proto.NewTaskServiceClient(conn),
		conn:   conn,
	}

	t.Run("SendAndReceiveResult", func(t *testing.T) {
		testTask := &models.Task{ID: 1, Operation: "+", Arg1: 10, Arg2: 20}
		result, _, err := agent.executeTask(context.Background(), testTask)
		assert.NoError(t, err, "Task execution error")
		assert.Equal(t, 30.0, result, "Incrorrect result")

		err = agent.sendResult(context.Background(), 1, result, "")
		assert.NoError(t, err, "Error sending result")

		select {
		case req := <-testServer.received:
			assert.Equal(t, float64(30), req.Result, "The server received an incorrect result")
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout waiting for result")
		}
	})
}

type mockTaskServiceClient struct {
	proto.TaskServiceClient
	receiveTaskError bool
}

func (m *mockTaskServiceClient) ReceiveTask(ctx context.Context, req *proto.GetTaskRequest, opts ...grpc.CallOption) (*proto.Task, error) {
	if m.receiveTaskError {
		return nil, fmt.Errorf("mock receive error")
	}
	return &proto.Task{}, nil
}

func TestGetTask_Error(t *testing.T) {
	agent := &GRPCAgent{
		client: &mockTaskServiceClient{receiveTaskError: true},
	}
	_, err := agent.getTask(context.Background())
	assert.Error(t, err, "Error expected but got nil")
}

func (m *mockTaskServiceClient) SubmitTaskResult(ctx context.Context, req *proto.SubmitTaskResultRequest, opts ...grpc.CallOption) (*proto.SubmitTaskResultResponse, error) {
	return nil, fmt.Errorf("mock submit error")
}

func TestSendResult_Error(t *testing.T) {
	agent := &GRPCAgent{
		client: &mockTaskServiceClient{},
	}
	err := agent.sendResult(context.Background(), 1, 10, "")
	assert.Error(t, err, "Error sending result expected")
}

func TestExecuteTask_NilTask(t *testing.T) {
	agent := &GRPCAgent{}
	_, _, err := agent.executeTask(context.Background(), nil)
	assert.Error(t, err, "Error expected for nil task")
}

func TestExecuteTask_ZeroID(t *testing.T) {
	agent := &GRPCAgent{}
	task := &models.Task{ID: 0}
	_, _, err := agent.executeTask(context.Background(), task)
	assert.Error(t, err, "Error expected for task with ID=0")
}

func TestExecuteTask_EmptyOperation(t *testing.T) {
	agent := &GRPCAgent{}
	task := &models.Task{ID: 1, Operation: ""}
	_, _, err := agent.executeTask(context.Background(), task)
	assert.Error(t, err, "An empty operation error was expected.")
}

func TestExecuteTask_DivisionByZero(t *testing.T) {
	agent := &GRPCAgent{}
	task := &models.Task{ID: 1, Operation: "/", Arg1: 5, Arg2: 0}
	result, errMsg, err := agent.executeTask(context.Background(), task)
	assert.NoError(t, err, "Expected nil, but got error")
	assert.Equal(t, models.ErrorDivisionByZero.Error(), errMsg, "Expected division by zero error message")
	assert.Equal(t, 0.0, result, "Result should be 0")
}

func TestExecuteTask_OperationsWithDelays(t *testing.T) {
	agent := &GRPCAgent{}
	tests := []struct {
		name           string
		operation      string
		arg1, arg2     float64
		delayKey       string
		expectedResult float64
	}{
		{"Addition", "+", 5, 3, "duration.TIME_ADDITION_MS", 8},
		{"Subtraction", "-", 10, 7, "duration.TIME_SUBTRACTION_MS", 3},
		{"Multiplication", "*", 4, 5, "duration.TIME_MULTIPLICATIONS_MS", 20},
		{"Division", "/", 10, 2, "duration.TIME_DIVISIONS_MS", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(tt.delayKey, 100)
			task := &models.Task{
				ID:        1,
				Arg1:      tt.arg1,
				Arg2:      tt.arg2,
				Operation: tt.operation,
			}
			start := time.Now()
			result, errMsg, err := agent.executeTask(context.Background(), task)
			elapsed := time.Since(start)

			assert.NoError(t, err)
			assert.Empty(t, errMsg)
			assert.Equal(t, tt.expectedResult, result)
			assert.InDelta(t, 100, elapsed.Milliseconds(), 50)
		})
	}
}

func TestNewGRPCAgent_UsesDefaults(t *testing.T) {
	viper.Reset()

	agent, err := newGRPCAgent()
	assert.NoError(t, err)
	assert.NotNil(t, agent)

	target := agent.conn.Target()
	assert.Equal(t, "localhost:50051", target)
}

func TestSendResult_WithErrorMessage(t *testing.T) {
	srv, testServer := startTestServer()
	defer srv.Stop()

	conn, _ := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	agent := &GRPCAgent{
		client: proto.NewTaskServiceClient(conn),
		conn:   conn,
	}

	err := agent.sendResult(context.Background(), 1, 0, "division by zero")
	assert.NoError(t, err)

	select {
	case req := <-testServer.received:
		assert.Equal(t, "division by zero", req.ErrorMessage)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for result")
	}
}

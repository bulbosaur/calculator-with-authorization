package agent

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
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

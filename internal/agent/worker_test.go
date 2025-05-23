package agent

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestWorker(t *testing.T) {
	srv, testServer := startTestServer()
	defer srv.Stop()

	testServer.tasks[1] = &proto.Task{
		Id:           1,
		ExpressionId: 1,
		Arg1:         15,
		Arg2:         5,
		Operation:    "-",
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
		Client: proto.NewTaskServiceClient(conn),
		Conn:   conn,
	}

	Workers = 1
	go agent.Worker(1)

	select {
	case req := <-testServer.received:
		assert.Equal(t, float64(10), req.Result, "Incorrect result")
		assert.Equal(t, int32(1), req.TaskId, "incorrect task ID")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for result")
	}
}

func TestExecutionTime(t *testing.T) {
	agent := &GRPCAgent{}
	task := &models.Task{
		ID:        1,
		Arg1:      2,
		Arg2:      3,
		Operation: "+",
	}

	start := time.Now()
	result, errMsg, err := agent.executeTask(context.Background(), task)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Empty(t, errMsg)
	assert.InDelta(t, viper.GetInt("duration.TIME_ADDITION_MS"), elapsed.Milliseconds(), 50)
	assert.Equal(t, float64(5), result)
}

func TestWorkerErrorHandling(t *testing.T) {
	srv, testServer := startTestServer()
	defer srv.Stop()

	testServer.tasks[1] = &proto.Task{
		Id:        1,
		Arg1:      0,
		Arg2:      0,
		Operation: "invalid_operation",
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
		Client: proto.NewTaskServiceClient(conn),
		Conn:   conn,
	}

	logChan := make(chan string, 2)
	log.SetOutput(&logWriter{logs: logChan})
	defer log.SetOutput(os.Stderr)

	Workers = 1
	go agent.Worker(1)

	var logMsg string
	for i := 0; i < 2; i++ {
		select {
		case msg := <-logChan:
			logMsg += msg
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout waiting for log message")
		}
	}

	assert.Contains(t, logMsg, "execution error task ID-1", "No execution error log message")
}

type logWriter struct {
	logs chan string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.logs <- string(p)
	return len(p), nil
}

func TestWorkerNilTaskHandling(t *testing.T) {
	srv, testServer := startTestServer()
	defer srv.Stop()

	testServer.tasks[1] = &proto.Task{
		Id: 0,
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
		Client: proto.NewTaskServiceClient(conn),
		Conn:   conn,
	}

	logChan := make(chan string, 1)
	log.SetOutput(&logWriter{logs: logChan})
	defer log.SetOutput(os.Stderr)

	Workers = 1
	go agent.Worker(1)

	select {
	case msg := <-logChan:
		assert.Contains(t, msg, "invalid task: task is nil or has ID 0")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for log message")
	}
}

func TestWorker_RetryOnFailure(t *testing.T) {
	mockClient := &mockTaskServiceClient{receiveTaskError: true}
	agent := &GRPCAgent{Client: mockClient}

	logChan := make(chan string, 1)
	log.SetOutput(&logWriter{logs: logChan})
	defer log.SetOutput(os.Stderr)

	Workers = 1
	go agent.Worker(1)

	select {
	case msg := <-logChan:
		assert.Contains(t, msg, "task receiving error")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for retry log")
	}
}

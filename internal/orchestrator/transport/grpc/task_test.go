package orchestrator_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	orchestrator "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/service"
	orchestratorGrpc "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/grpc"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "modernc.org/sqlite"
)

type testServer struct {
	proto.UnimplementedTaskServiceServer
	exprRepo *repository.ExpressionModel
	server   *grpc.Server
	port     string
	db       *sql.DB
}

func setupTestServer(t *testing.T) *testServer {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS expressions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            expression TEXT NOT NULL,
            status TEXT NOT NULL,
            result FLOAT64 DEFAULT 0,
            error_message TEXT DEFAULT ""
        );
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            expressionID INTEGER NOT NULL,
            arg1 TEXT NOT NULL,
            arg2 TEXT NOT NULL,
            prev_task_id1 INTEGER DEFAULT 0,
            prev_task_id2 INTEGER DEFAULT 0,
            operation TEXT NOT NULL,
            status TEXT,
            result FLOAT,
            error_message TEXT DEFAULT ""
        );
    `)
	require.NoError(t, err)

	exprRepo := &repository.ExpressionModel{DB: db}

	port := "50051"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	require.NoError(t, err)

	server := grpc.NewServer()
	proto.RegisterTaskServiceServer(server, &orchestratorGrpc.TaskServer{
		ExprRepo: exprRepo,
	})

	go func() {
		log.Printf("gRPC test server listening on port %s", port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return &testServer{
		exprRepo: exprRepo,
		server:   server,
		port:     port,
		db:       db,
	}
}

func (ts *testServer) teardown(t *testing.T) {
	ts.server.Stop()
	err := ts.db.Close()
	require.NoError(t, err)
}

func TestReceiveTask(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.teardown(t)

	exprID, err := ts.exprRepo.Insert("3+4", 1)
	require.NoError(t, err)

	err = orchestrator.Calc("3+4", exprID, ts.exprRepo)
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", ts.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := proto.NewTaskServiceClient(conn)

	resp, err := client.ReceiveTask(context.Background(), &proto.GetTaskRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int32(exprID), resp.ExpressionId)
	assert.Equal(t, float64(3), resp.Arg1)
	assert.Equal(t, float64(4), resp.Arg2)
	assert.Equal(t, "+", resp.Operation)

	task, err := ts.exprRepo.GetTaskByID(int(resp.Id))
	require.NoError(t, err)
	assert.Equal(t, models.StatusInProcess, task.Status)
}

package agent

// import (
// 	"context"
// 	"net"
// 	"testing"
// 	"time"

// 	"github.com/bulbosaur/calculator-with-authorization/proto"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func TestWorkerIntegration(t *testing.T) {
// 	srv, testServer := StartTestServer()
// 	defer srv.Stop()

// 	testServer.tasks[1] = &proto.Task{
// 		Id:        1,
// 		Arg1:      15,
// 		Arg2:      25,
// 		Operation: "+",
// 	}

// 	conn, err := grpc.NewClient(
// 		"bufnet",
// 		context.Background(),
// 		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
// 			return lis.Dial()
// 		}),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	assert.NoError(t, err)

// 	agent := &GRPCAgent{
// 		client: proto.NewTaskServiceClient(conn),
// 		conn:   conn,
// 	}
// 	Workers = 1

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	go agent.worker(1)

// 	select {
// 	case req := <-testServer.received:
// 		assert.Equal(t, float64(40), req.Result)
// 	case <-ctx.Done():
// 		t.Fatal("Timeout waiting for worker processing")
// 	}
// }

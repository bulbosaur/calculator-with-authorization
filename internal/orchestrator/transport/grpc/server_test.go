package orchestrator_test

import (
	"net"
	"strings"
	"testing"

	orchestratorGRPC "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/grpc"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGRPCOrchestrator_Failure(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	defer listener.Close()

	addr := listener.Addr().String()
	hostPort := strings.Split(addr, ":")
	require.Len(t, hostPort, 2)
	port := hostPort[1]

	viper.Set("server.GRPC_HOST", "localhost")
	viper.Set("server.GRPC_PORT", port)
	defer viper.Reset()

	exprRepo := &repository.ExpressionModel{DB: nil}

	err = orchestratorGRPC.RunGRPCOrchestrator(exprRepo)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to listen")
}

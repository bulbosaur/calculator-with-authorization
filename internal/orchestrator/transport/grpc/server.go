package orchestrator

import (
	"fmt"
	"log"
	"net"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// RunGRPCOrchestrator запускает gRPC сервер оркестратора
func RunGRPCOrchestrator(exprRepo models.ExpressionRepository) error {
	host := viper.GetString("server.GRPC_HOST")
	port := viper.GetString("server.GRPC_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024),
		grpc.MaxSendMsgSize(10*1024*1024),
	)
	proto.RegisterTaskServiceServer(s, newTaskServer(exprRepo))
	log.Printf("gRPC server listening on %s", addr)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

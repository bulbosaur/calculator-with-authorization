package orchestrator

import (
	"fmt"
	"log"
	"net"

	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/bulbosaur/calculator-with-authorization/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// RunGRPCServer запускает gRPC сервер оркестратора
func RunGRPCOrchestrator(exprRepo *repository.ExpressionModel) {
	host := viper.GetString("server.GRPC_HOST")
	port := viper.GetString("server.GRPC_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, NewTaskServer(exprRepo))

	log.Printf("gRPC server listening on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

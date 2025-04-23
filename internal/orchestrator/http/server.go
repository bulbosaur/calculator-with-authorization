package orchestrator

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// RunOrchestrator запускает сервер оркестратора
func RunOrchestrator(exprRepo *repository.ExpressionModel) {
	go func() {
		runGRPCServer(exprRepo)
	}()

	host := viper.GetString("server.ORC_HOST")
	port := viper.GetString("server.ORC_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/api/v1/calculate", regHandler(exprRepo)).Methods("POST")
	router.HandleFunc("/api/v1/expressions", listHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", resultHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/coffee", CoffeeHandler)

	log.Printf("Orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("Orchestrator server error:", err)
	}
}

func runGRPCServer(exprRepo *repository.ExpressionModel) {
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

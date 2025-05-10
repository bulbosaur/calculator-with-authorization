package main

import (
	"log"

	config "github.com/bulbosaur/calculator-with-authorization/config"
	orchestratorGRPC "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/grpc"
	orchestratorHTTP "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/spf13/viper"

	_ "modernc.org/sqlite"
)

func main() {
	log.Println("Starting server...")

	config.Init()

	db, err := repository.InitDB(viper.GetString("DATABASE_PATH"))
	if err != nil {
		log.Fatalf("failed to init DB; %v", err)
	}

	ExprRepo := repository.NewExpressionModel(db)

	defer db.Close()

	go orchestratorHTTP.RunHTTPOrchestrator(ExprRepo)
	err = orchestratorGRPC.RunGRPCOrchestrator(ExprRepo)

	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

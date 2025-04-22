package main

import (
	"log"

	config "github.com/bulbosaur/calculator-with-authorization/config"
	orchestrator "github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/http"
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

	orchestrator.RunOrchestrator(ExprRepo)
}

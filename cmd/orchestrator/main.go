package main

import (
	"log"

	config "github.com/bulbosaur/web-calculator-golang/config"
	orchestrator "github.com/bulbosaur/web-calculator-golang/internal/orchestrator/http"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
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

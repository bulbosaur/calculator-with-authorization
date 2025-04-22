package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// RunOrchestrator запускает сервер оркестратора
func RunOrchestrator(exprRepo *repository.ExpressionModel) {
	host := viper.GetString("server.ORC_HOST")
	port := viper.GetString("server.ORC_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/api/v1/calculate", regHandler(exprRepo)).Methods("POST")
	router.HandleFunc("/api/v1/expressions", listHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", resultHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/internal/task", taskHandler(exprRepo))
	router.HandleFunc("/internal/task/{id}", taskResultHandler(exprRepo))
	router.HandleFunc("/coffee", CoffeeHandler)

	log.Printf("Orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("Orchestrator server error:", err)
	}
}

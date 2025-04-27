package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// RunHTTPOrchestrator запускает http сервер оркестратора
func RunHTTPOrchestrator(exprRepo *repository.ExpressionModel) {

	host := viper.GetString("server.HTTP_HOST")
	port := viper.GetString("server.HTTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/api/v1/calculate", handlers.RegHandler(exprRepo)).Methods("POST")
	router.HandleFunc("/api/v1/expressions", handlers.ListHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", handlers.ResultHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/coffee", handlers.CoffeeHandler)

	log.Printf("HTTP orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("HTTP orchestrator server error:", err)
	}
}

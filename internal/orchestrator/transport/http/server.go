package orchestrator

import (
	"fmt"
	"log"
	"net/http"

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

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/api/v1/calculate", regHandler(exprRepo)).Methods("POST")
	router.HandleFunc("/api/v1/expressions", listHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", resultHandler(exprRepo)).Methods("GET")
	router.HandleFunc("/coffee", CoffeeHandler)

	log.Printf("HTTP orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("HTTP orchestrator server error:", err)
	}
}

package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/config"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/middlewares"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// RunHTTPOrchestrator запускает http сервер оркестратора
func RunHTTPOrchestrator(exprRepo *repository.ExpressionModel, cfg *config.JWTConfig) {

	host := viper.GetString("server.HTTP_HOST")
	port := viper.GetString("server.HTTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/login", handlers.LoginHandler(exprRepo, cfg)).Methods("POST")
	router.HandleFunc("/api/v1/register", handlers.Register(exprRepo)).Methods("POST")

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middlewares.AuthMiddleware(cfg))

	protected.HandleFunc("/calculate", handlers.RegHandler(exprRepo)).Methods("POST")
	protected.HandleFunc("/expressions", handlers.ListHandler(exprRepo)).Methods("GET")
	protected.HandleFunc("/expressions/{id}", handlers.ResultHandler(exprRepo)).Methods("GET")
	log.Printf("HTTP orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("HTTP orchestrator server error:", err)
	}
}

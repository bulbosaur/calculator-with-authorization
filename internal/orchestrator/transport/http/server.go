package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/middlewares"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// RunHTTPOrchestrator запускает http сервер оркестратора
func RunHTTPOrchestrator(exprRepo *repository.ExpressionModel) {

	host := viper.GetString("server.HTTP_HOST")
	port := viper.GetString("server.HTTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	Service := auth.NewService(viper.GetString("jwt.secret_key"), viper.GetDuration("jwt.token_duration"))

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))),
	)

	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	router.HandleFunc("/login", handlers.LoginPageHandler).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterPageHandler).Methods("GET")

	router.HandleFunc("/coffee", handlers.CoffeeHandler).Methods("GET")

	router.HandleFunc("/api/v1/login", handlers.LoginHandler(Service, exprRepo)).Methods("POST")
	router.HandleFunc("/api/v1/register", handlers.Register(Service, exprRepo)).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(middlewares.AuthMiddleware(Service))

	protected.HandleFunc("/calculator", handlers.CalcPageHandler).Methods("GET")

	protected.HandleFunc("/api/v1/calculate", handlers.RegHandler(exprRepo)).Methods("POST")
	protected.HandleFunc("/api/v1/expressions", handlers.ListHandler(Service, exprRepo)).Methods("GET")
	protected.HandleFunc("/api/v1/expressions/{id}", handlers.ResultHandler(Service, exprRepo)).Methods("GET")

	log.Printf("HTTP orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("HTTP orchestrator server error:", err)
	}
}

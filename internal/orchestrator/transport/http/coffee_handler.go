package orchestrator

import (
	"fmt"
	"net/http"
)

// CoffeeHandler - сердце программы, "Я чайник"
func CoffeeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "I'm a teapot")
}

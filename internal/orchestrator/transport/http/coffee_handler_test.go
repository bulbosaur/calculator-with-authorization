package orchestrator

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoffeeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/coffee", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CoffeeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTeapot {
		t.Errorf("CoffeHandler returned %v, but want %v",
			status, http.StatusTeapot)
	}

	expected := "I'm a teapot"
	if rr.Body.String() != expected {
		t.Errorf("CoffeHandler returned %v, but want %v",
			rr.Body.String(), expected)
	}
}

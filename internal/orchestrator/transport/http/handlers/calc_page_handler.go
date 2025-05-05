package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func CalcPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlpath := filepath.Join("web", "templates", "calc.html")
	htmlContent, err := os.ReadFile(htmlpath)
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

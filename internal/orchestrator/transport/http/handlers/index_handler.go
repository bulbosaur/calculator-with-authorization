package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// IndexHandler - главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join("web", "templates", "index.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

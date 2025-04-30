package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// RegisterPageHandler - обработчик страницы регистрации
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join("web", "templates", "register.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		http.Error(
			w,
			"Registration page is unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

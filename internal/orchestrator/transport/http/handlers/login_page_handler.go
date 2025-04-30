package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// LoginPageHandler - обработчик страницы авторизации
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join("web", "templates", "login.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		http.Error(
			w,
			"Authorization page is unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

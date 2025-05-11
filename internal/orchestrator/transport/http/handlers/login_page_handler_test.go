package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/orchestrator/transport/http/handlers"
	"github.com/stretchr/testify/assert"
)

func TestLoginPageHandler_Success(t *testing.T) {
	tmpDir := t.TempDir()

	templateDir := filepath.Join(tmpDir, "web", "templates")
	err := os.MkdirAll(templateDir, os.ModePerm)
	assert.NoError(t, err)

	htmlContent := "<html><body>Login Page</body></html>"
	htmlPath := filepath.Join(templateDir, "login.html")
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	assert.NoError(t, err)

	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	handlers.LoginPageHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<html>")
}

func TestLoginPageHandler_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	handlers.LoginPageHandler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization page is unavailable")
}

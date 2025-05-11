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

func TestCalcPageHandler(t *testing.T) {
	tmpDir := t.TempDir()

	templateDir := filepath.Join(tmpDir, "web", "templates")
	err := os.MkdirAll(templateDir, os.ModePerm)
	assert.NoError(t, err)

	htmlContent := "<html><body>Test Calculator Page</body></html>"
	htmlPath := filepath.Join(templateDir, "calc.html")
	err = os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	assert.NoError(t, err)

	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/calc", nil)
	w := httptest.NewRecorder()
	handlers.CalcPageHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<html>")
}

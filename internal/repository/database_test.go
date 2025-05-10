package repository_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var repo *repository.ExpressionModel

func setupTestDB(t *testing.T) func() {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	var err error
	db, err = repository.InitDB(dbPath)
	assert.NoError(t, err)

	repo = repository.NewExpressionModel(db)

	return func() {
		db.Close()
		os.RemoveAll(tempDir)
	}
}

func TestInsertAndGetExpression(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, err := repo.Insert("2 + 3", 1)
	assert.NoError(t, err)
	assert.Greater(t, exprID, 0)

	expr, err := repo.GetExpression(exprID)
	assert.NoError(t, err)
	assert.Equal(t, "2 + 3", expr.Expression)
	assert.Equal(t, models.StatusWait, expr.Status)
	assert.Equal(t, 0.0, expr.Result)
}

func TestUpdateExpressionStatus(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("5 * 5", 1)

	repo.UpdateStatus(exprID, models.StatusInProcess)
	expr, _ := repo.GetExpression(exprID)
	assert.Equal(t, models.StatusInProcess, expr.Status)
}

func TestUpdateExpressionResult(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	exprID, _ := repo.Insert("10 / 2", 1)

	err := repo.UpdateExpressionResult(exprID, 5.0, "")
	assert.NoError(t, err)

	expr, _ := repo.GetExpression(exprID)
	assert.Equal(t, models.StatusResolved, expr.Status)
	assert.Equal(t, 5.0, expr.Result)
}

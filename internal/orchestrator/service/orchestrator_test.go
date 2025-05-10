package orchestrator

import (
	"database/sql"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
	_ "modernc.org/sqlite"
)

func TestCalc(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expressions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        expression TEXT NOT NULL,
        status TEXT NOT NULL,
        result FLOAT64 DEFAULT 0,
        error_message TEXT DEFAULT ""
    );`)
	if err != nil {
		t.Fatalf("Failed to create expressions table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        expressionID INTEGER NOT NULL,
        arg1 TEXT NOT NULL,
        arg2 TEXT NOT NULL,
        prev_task_id1 INTEGER DEFAULT 0,
        prev_task_id2 INTEGER DEFAULT 0,
        operation TEXT NOT NULL,
        status TEXT,
        result FLOAT,
        error_message TEXT DEFAULT ""
    );`)
	if err != nil {
		t.Fatalf("Failed to create tasks table: %v", err)
	}

	repo := &repository.ExpressionModel{DB: db}

	tests := []struct {
		name        string
		expression  string
		expectError error
	}{
		{
			name:        "EmptyExpression",
			expression:  "",
			expectError: models.ErrorEmptyExpression,
		},
		{
			name:        "InvalidCharacter",
			expression:  "2@2",
			expectError: models.ErrorInvalidCharacter,
		},
		{
			name:        "UnclosedBracket",
			expression:  "(3+4",
			expectError: models.ErrorUnclosedBracket,
		},
		{
			name:        "ValidExpression",
			expression:  "3+4*2",
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Calc(tt.expression, 1, repo)
			if (err != nil && tt.expectError == nil) || (err == nil && tt.expectError != nil) {
				t.Errorf("Unexpected error: got %v, want %v", err, tt.expectError)
			} else if err != nil && tt.expectError != nil && err.Error() != tt.expectError.Error() {
				t.Errorf("Error mismatch: got %v, want %v", err, tt.expectError)
			}
		})
	}
}

package orchestrator

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/bulbosaur/calculator-with-authorization/internal/repository"
)

func tokensEqual(a, b []models.Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestToReversePolishNotation(t *testing.T) {
	tests := []struct {
		input    []models.Token
		expected []models.Token
		err      error
	}{
		{
			input: []models.Token{
				{Value: "3", IsNumber: true},
				{Value: "+", IsNumber: false},
				{Value: "4", IsNumber: true},
			},
			expected: []models.Token{
				{Value: "3", IsNumber: true},
				{Value: "4", IsNumber: true},
				{Value: "+", IsNumber: false},
			},
			err: nil,
		},
		{
			input: []models.Token{
				{Value: "(", IsNumber: false},
				{Value: "3", IsNumber: true},
				{Value: "+", IsNumber: false},
				{Value: "4", IsNumber: true},
				{Value: ")", IsNumber: false},
				{Value: "*", IsNumber: false},
				{Value: "5", IsNumber: true},
			},
			expected: []models.Token{
				{Value: "3", IsNumber: true},
				{Value: "4", IsNumber: true},
				{Value: "+", IsNumber: false},
				{Value: "5", IsNumber: true},
				{Value: "*", IsNumber: false},
			},
			err: nil,
		},
		{
			input: []models.Token{
				{Value: "4", IsNumber: true},
				{Value: "!", IsNumber: false},
				{Value: "2", IsNumber: true},
			},
			expected: nil,
			err:      models.ErrorInvalidInput,
		},
	}

	for i, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := toReversePolishNotation(tt.input)
			if err != tt.err || !tokensEqual(result, tt.expected) {
				t.Errorf("Test case %d failed: got %v, want %v, err %v, wantErr %v", i, result, tt.expected, err, tt.err)
			}
		})
	}
}

func TestParseRPN(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	db.Exec(`
        CREATE TABLE tasks (
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
        )
    `)
	repo := &repository.ExpressionModel{DB: db}

	tests := []struct {
		name   string
		tokens []models.Token
		err    error
	}{
		{
			name: "SimpleAddition",
			tokens: []models.Token{
				{Value: "3", IsNumber: true},
				{Value: "4", IsNumber: true},
				{Value: "+", IsNumber: false},
			},
			err: nil,
		},
		{
			name: "NotEnoughOperands",
			tokens: []models.Token{
				{Value: "3", IsNumber: true},
				{Value: "+", IsNumber: false},
			},
			err: fmt.Errorf("not enough operands for operation +"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseRPN(tt.tokens, 1, repo)
			if (err == nil && tt.err != nil) || (err != nil && tt.err == nil) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			} else if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Error mismatch: got %v, want %v", err, tt.err)
			}
		})
	}
}

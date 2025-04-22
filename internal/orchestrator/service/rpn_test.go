package orchestrator

import (
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
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

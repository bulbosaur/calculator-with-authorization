package orchestrator

import (
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

func equalTokens(a, b []models.Token) bool {
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

func TestTokenize(t *testing.T) {
	cases := []struct {
		expression string
		wantOutput []models.Token
		wantError  error
	}{
		{
			expression: "2+2",
			wantOutput: []models.Token{{"2", true}, {"+", false}, {"2", true}},
			wantError:  nil,
		},
		{
			expression: "12 + 34",
			wantOutput: []models.Token{{"12", true}, {"+", false}, {"34", true}},
			wantError:  nil,
		},
		{
			expression: "(5-3)*2",
			wantOutput: []models.Token{
				{"(", false},
				{"5", true},
				{"-", false},
				{"3", true},
				{")", false},
				{"*", false},
				{"2", true},
			},
			wantError: nil,
		},
		{
			expression: "42",
			wantOutput: []models.Token{{"42", true}},
			wantError:  nil,
		},
		{
			expression: "(2 - 5)(3 + 4)",
			wantOutput: []models.Token{
				{"(", false},
				{"2", true},
				{"-", false},
				{"5", true},
				{")", false},
				{"*", false},
				{"(", false},
				{"3", true},
				{"+", false},
				{"4", true},
				{")", false},
			},
			wantError: nil,
		},
		{
			expression: "",
			wantOutput: []models.Token{},
			wantError:  nil,
		},
		{
			expression: "2 @ 3",
			wantOutput: nil,
			wantError:  models.ErrorInvalidCharacter,
		},
		{
			expression: "32 3",
			wantOutput: nil,
			wantError:  models.ErrorMissingOperand,
		},
		{
			expression: "2.5 + 3",
			wantOutput: nil,
			wantError:  models.ErrorInvalidCharacter,
		},
		{
			expression: "()",
			wantOutput: nil,
			wantError:  models.ErrorEmptyBrackets,
		},
		{
			expression: "a + b",
			wantOutput: nil,
			wantError:  models.ErrorInvalidCharacter,
		},
		{
			expression: "45(4 -8",
			wantOutput: nil,
			wantError:  models.ErrorUnclosedBracket,
		},
	}
	for _, tc := range cases {
		t.Run(tc.expression, func(t *testing.T) {
			got, err := tokenize(tc.expression)
			if err != tc.wantError {
				t.Errorf("Tokenize(%v) error = %v, wantErr %v", tc.expression, err, tc.wantError)
				return
			}

			if tc.wantError != nil {
				return
			}

			if !equalTokens(got, tc.wantOutput) {
				t.Errorf("Calc(%v) = %v, want %v", tc.expression, got, tc.wantOutput)
			}
		})
	}
}

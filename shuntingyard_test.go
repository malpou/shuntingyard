package shuntingyard

import (
	"math"
	"testing"
)

// Helper function to compare floats with tolerance
func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// TestScan tests the tokenization of mathematical expressions
func TestScan(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "simple addition",
			input:    "2 + 3",
			expected: []string{"2", "+", "3"},
			wantErr:  false,
		},
		{
			name:     "no spaces",
			input:    "1+2",
			expected: []string{"1", "+", "2"},
			wantErr:  false,
		},
		{
			name:     "mixed spacing",
			input:    "1 + 2+3",
			expected: []string{"1", "+", "2", "+", "3"},
			wantErr:  false,
		},
		{
			name:     "floating point numbers",
			input:    "1.5 + 2.5",
			expected: []string{"1.5", "+", "2.5"},
			wantErr:  false,
		},
		{
			name:     "complex expression",
			input:    "10.5 / 2 + 3.5",
			expected: []string{"10.5", "/", "2", "+", "3.5"},
			wantErr:  false,
		},
		{
			name:     "with parentheses",
			input:    "(2 + 3) * 4",
			expected: []string{"(", "2", "+", "3", ")", "*", "4"},
			wantErr:  false,
		},
		{
			name:     "multi-digit numbers",
			input:    "100 / 2 - 3 * 4 + 5",
			expected: []string{"100", "/", "2", "-", "3", "*", "4", "+", "5"},
			wantErr:  false,
		},
		{
			name:    "invalid character",
			input:   "2 + 3a",
			wantErr: true,
		},
		{
			name:    "empty expression",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Scan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Scan() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Scan() unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Scan() got %d tokens, expected %d", len(result), len(tt.expected))
				return
			}

			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("Scan() token[%d] = %s, expected %s", i, token, tt.expected[i])
				}
			}
		})
	}
}

// TestParse tests the conversion from infix to postfix notation
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "simple addition",
			input:    []string{"2", "+", "3"},
			expected: []string{"2", "3", "+"},
			wantErr:  false,
		},
		{
			name:     "operator precedence",
			input:    []string{"2", "+", "3", "*", "4"},
			expected: []string{"2", "3", "4", "*", "+"},
			wantErr:  false,
		},
		{
			name:     "division and subtraction precedence",
			input:    []string{"10", "-", "6", "/", "2"},
			expected: []string{"10", "6", "2", "/", "-"},
			wantErr:  false,
		},
		{
			name:     "parentheses override precedence",
			input:    []string{"(", "2", "+", "3", ")", "*", "4"},
			expected: []string{"2", "3", "+", "4", "*"},
			wantErr:  false,
		},
		{
			name:     "nested parentheses",
			input:    []string{"(", "(", "2", "+", "3", ")", "*", "4", ")", "-", "5"},
			expected: []string{"2", "3", "+", "4", "*", "5", "-"},
			wantErr:  false,
		},
		{
			name:     "multiple parentheses groups",
			input:    []string{"(", "2", "+", "3", ")", "*", "(", "4", "+", "1", ")"},
			expected: []string{"2", "3", "+", "4", "1", "+", "*"},
			wantErr:  false,
		},
		{
			name:     "left associativity",
			input:    []string{"1", "+", "2", "+", "3", "+", "4", "+", "5"},
			expected: []string{"1", "2", "+", "3", "+", "4", "+", "5", "+"},
			wantErr:  false,
		},
		{
			name:    "mismatched parentheses - extra right",
			input:   []string{"2", "+", "3", ")"},
			wantErr: true,
		},
		{
			name:    "mismatched parentheses - extra left",
			input:   []string{"(", "2", "+", "3"},
			wantErr: true,
		},
		{
			name:    "empty token list",
			input:   []string{},
			wantErr: true,
		},
		{
			name:    "invalid token",
			input:   []string{"2", "+", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Parse() got %d tokens, expected %d", len(result), len(tt.expected))
				t.Errorf("Got: %v, Expected: %v", result, tt.expected)
				return
			}

			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("Parse() token[%d] = %s, expected %s", i, token, tt.expected[i])
				}
			}
		})
	}
}

// TestEvaluate tests the evaluation of postfix expressions
func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected float64
		wantErr  bool
	}{
		{
			name:     "simple addition",
			input:    []string{"2", "3", "+"},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:     "simple subtraction",
			input:    []string{"5", "3", "-"},
			expected: 2.0,
			wantErr:  false,
		},
		{
			name:     "simple multiplication",
			input:    []string{"4", "5", "*"},
			expected: 20.0,
			wantErr:  false,
		},
		{
			name:     "simple division",
			input:    []string{"10", "2", "/"},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:     "floating point division",
			input:    []string{"10", "4", "/"},
			expected: 2.5,
			wantErr:  false,
		},
		{
			name:     "complex expression",
			input:    []string{"2", "3", "4", "*", "+"},
			expected: 14.0,
			wantErr:  false,
		},
		{
			name:     "floating point expression",
			input:    []string{"10.5", "2", "/", "3.5", "+"},
			expected: 8.75,
			wantErr:  false,
		},
		{
			name:     "division resulting in repeating decimal",
			input:    []string{"10", "3", "/"},
			expected: 3.333333333333333,
			wantErr:  false,
		},
		{
			name:    "division by zero",
			input:   []string{"10", "0", "/"},
			wantErr: true,
		},
		{
			name:    "insufficient operands",
			input:   []string{"2", "+"},
			wantErr: true,
		},
		{
			name:    "too many operands",
			input:   []string{"2", "3", "4", "+"},
			wantErr: true,
		},
		{
			name:    "empty expression",
			input:   []string{},
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   []string{"abc", "2", "+"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Evaluate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Evaluate() unexpected error: %v", err)
				return
			}

			if !almostEqual(result, tt.expected, 0.0000001) {
				t.Errorf("Evaluate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestFullPipeline tests complete expression evaluation from scan to result
func TestFullPipeline(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
		wantErr    bool
	}{
		// Basic arithmetic
		{name: "basic addition", expression: "2 + 3", expected: 5.0},
		{name: "basic subtraction", expression: "5 - 3", expected: 2.0},
		{name: "basic multiplication", expression: "4 * 5", expected: 20.0},
		{name: "basic division", expression: "10 / 2", expected: 5.0},

		// Order of operations
		{name: "mult before add", expression: "2 + 3 * 4", expected: 14.0},
		{name: "div before sub", expression: "10 - 6 / 2", expected: 7.0},
		{name: "complex precedence", expression: "2 + 3 * 4 - 5", expected: 9.0},

		// Parentheses
		{name: "simple parens", expression: "(2 + 3) * 4", expected: 20.0},
		{name: "nested parens", expression: "((2 + 3) * 4) - 5", expected: 15.0},
		{name: "multiple parens", expression: "(2 + 3) * (4 + 1)", expected: 25.0},

		// Floating-point numbers
		{name: "float division", expression: "10 / 4", expected: 2.5},
		{name: "float addition", expression: "1.5 + 2.5", expected: 4.0},
		{name: "float complex", expression: "10.5 / 2 + 3.5", expected: 8.75},
		{name: "repeating decimal", expression: "10 / 3", expected: 3.333333333333333},

		// Spacing variations
		{name: "no spaces", expression: "1+2", expected: 3.0},
		{name: "mixed spacing", expression: "1 + 2+3", expected: 6.0},

		// Complex expressions
		{name: "many additions", expression: "1 + 2 + 3 + 4 + 5", expected: 15.0},
		{name: "mixed operations", expression: "100 / 2 - 3 * 4 + 5", expected: 43.0},

		// Error cases
		{name: "division by zero", expression: "10 / 0", wantErr: true},
		{name: "invalid character", expression: "2 + a", wantErr: true},
		{name: "mismatched parens", expression: "(2 + 3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Scan
			tokens, err := Scan(tt.expression)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Scan() unexpected error: %v", err)
				}
				return
			}

			// Step 2: Parse
			postfix, err := Parse(tokens)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() unexpected error: %v", err)
				}
				return
			}

			// Step 3: Evaluate
			result, err := Evaluate(postfix)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Evaluate() unexpected error: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("Expected error but got result: %v", result)
				return
			}

			if !almostEqual(result, tt.expected, 0.0000001) {
				t.Errorf("%s = %v, expected %v", tt.expression, result, tt.expected)
			}
		})
	}
}

// BenchmarkFullPipeline benchmarks the complete evaluation pipeline
func BenchmarkFullPipeline(b *testing.B) {
	expression := "100 / 2 - 3 * 4 + 5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokens, _ := Scan(expression)
		postfix, _ := Parse(tokens)
		_, _ = Evaluate(postfix)
	}
}

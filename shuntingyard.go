// Package shuntingyard implements the Shunting Yard algorithm for parsing and evaluating
// mathematical expressions with full floating-point number support.
package shuntingyard

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Scan tokenizes a mathematical expression string into individual tokens.
// It supports floating-point numbers, operators (+, -, *, /), and parentheses.
// Expressions can have spaces or be continuous (e.g., "1 + 2" or "1+2").
//
// Returns a slice of tokens or an error if invalid characters are encountered.
func Scan(expression string) ([]string, error) {
	if expression == "" {
		return nil, fmt.Errorf("empty expression")
	}

	var tokens []string
	var currentNumber strings.Builder

	for i, ch := range expression {
		switch {
		case unicode.IsDigit(ch) || ch == '.':
			// Build multi-digit numbers and decimals
			currentNumber.WriteRune(ch)

		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '(' || ch == ')':
			// Flush any accumulated number before adding operator/parenthesis
			if currentNumber.Len() > 0 {
				tokens = append(tokens, currentNumber.String())
				currentNumber.Reset()
			}
			tokens = append(tokens, string(ch))

		case unicode.IsSpace(ch):
			// Spaces separate tokens, flush any accumulated number
			if currentNumber.Len() > 0 {
				tokens = append(tokens, currentNumber.String())
				currentNumber.Reset()
			}

		default:
			return nil, fmt.Errorf("invalid character '%c' at position %d", ch, i)
		}
	}

	// Don't forget the last number
	if currentNumber.Len() > 0 {
		tokens = append(tokens, currentNumber.String())
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no valid tokens found")
	}

	return tokens, nil
}

// Parse converts infix notation tokens to postfix notation (Reverse Polish Notation)
// using the Shunting Yard algorithm. It handles operator precedence and associativity:
// - Multiplication and division have higher precedence than addition and subtraction
// - Operators of the same precedence are left-associative
//
// Returns postfix tokens or an error for mismatched parentheses.
func Parse(tokens []string) ([]string, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty token list")
	}

	var output []string
	var operatorStack []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			// Pop operators with greater or equal precedence (left-associative)
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top == "(" {
					break
				}
				if precedence[top] < precedence[token] {
					break
				}
				// Pop operator to output
				output = append(output, top)
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)

		case "(":
			operatorStack = append(operatorStack, token)

		case ")":
			// Pop until we find the matching left parenthesis
			found := false
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]

				if top == "(" {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, fmt.Errorf("mismatched parentheses: unmatched ')'")
			}

		default:
			// Must be a number, validate it
			if _, err := strconv.ParseFloat(token, 64); err != nil {
				return nil, fmt.Errorf("invalid number: %s", token)
			}
			output = append(output, token)
		}
	}

	// Pop remaining operators
	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		if top == "(" {
			return nil, fmt.Errorf("mismatched parentheses: unmatched '('")
		}
		output = append(output, top)
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return output, nil
}

// Evaluate computes the result of a postfix (RPN) expression.
// It uses a stack-based algorithm to process operators and operands.
//
// Returns the computed float64 result or an error for invalid expressions or division by zero.
func Evaluate(postfixTokens []string) (float64, error) {
	if len(postfixTokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}

	var stack []float64

	for _, token := range postfixTokens {
		switch token {
		case "+", "-", "*", "/":
			// Need at least 2 operands
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression: insufficient operands for operator '%s'", token)
			}

			// Pop two operands (note: order matters for - and /)
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			}

			stack = append(stack, result)

		default:
			// Must be a number
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token)
			}
			stack = append(stack, num)
		}
	}

	// Should have exactly one value left
	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: too many operands")
	}

	return stack[0], nil
}

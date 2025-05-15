package parser

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const IntApproxTolerance = 1e-10

type Parser struct {
	variables map[string]float64
}

func NewParser() *Parser {
	return &Parser{
		variables: make(map[string]float64),
	}
}

func (p *Parser) Parse(input string) ([]float64, error) {
	results := make([]float64, 0)
	debug := os.Getenv("DEBUG") != ""

	for stmt := range strings.SplitSeq(input, ";") {
		trimedStmt := strings.TrimSpace(stmt)
		if trimedStmt == "" {
			continue
		}

		// Show each statement if DEBUG is set
		if debug {
			fmt.Printf("[debug] Input: '%s'\r\n", trimedStmt)
		}

		lexer, err := NewLexer(trimedStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating lexer: %w", err)
		}

		// Show tokens if DEBUG is set
		if debug {
			fmt.Printf("[debug] Tokens: %s\r\n", lexer)
		}

		expr, err := NewExpressionFromLexer(lexer)
		if err != nil {
			return nil, fmt.Errorf("error creating expression: %w", err)
		}

		// Show expression if DEBUG is set
		if debug {
			fmt.Printf("[debug] Expression: %s\r\n", expr)
		}

		// Handle variable assignment
		if expr.IsOPAssignment() {
			varName, rhs := expr.GetAssignment()
			if varName != "" && rhs != nil {
				val, err := rhs.Evaluate(p.variables)
				if err != nil {
					return nil, fmt.Errorf("error evaluating assignment: %w", err)
				}
				p.variables[varName] = val

				// Print dividing line for readability if DEBUG is set
				if debug {
					fmt.Printf("------------------------\r\n")
				}

				continue
			}
		}

		result, err := expr.Evaluate(p.variables)
		if err != nil {
			return nil, fmt.Errorf("error evaluating expression: %w", err)
		}

		// Check if the result is approximately an integer for display
		// This is to handle cases like 1.99999999999 to 2
		rounded := math.Round(result)
		if math.Abs(rounded-result) < IntApproxTolerance {
			result = rounded
		}

		// Show result if DEBUG is set
		if debug {
			fmt.Printf("[debug] Evaluated: %s\r\n", strconv.FormatFloat(result, 'f', -1, 64))
		}

		// Print dividing line for readability if DEBUG is set
		if debug {
			fmt.Printf("------------------------\r\n")
		}

		results = append(results, result)
	}

	return results, nil
}

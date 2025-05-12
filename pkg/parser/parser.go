package parser

import (
	"fmt"
	"strings"
)

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

	for stmt := range strings.SplitSeq(input, ";") {
		stmt := strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		lexer, err := NewLexer(stmt)
		if err != nil {
			return nil, fmt.Errorf("error creating lexer: %w", err)
		}

		expr, err := NewExpressionFromLexer(lexer)
		if err != nil {
			return nil, fmt.Errorf("error creating expression: %w", err)
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
				continue
			}
		}

		result, err := expr.Evaluate(p.variables)
		if err != nil {
			return nil, fmt.Errorf("error evaluating expression: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

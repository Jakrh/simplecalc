package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"simplecalc/pkg/parser"
	"simplecalc/pkg/terminal"
)

func main() {
	t, err := terminal.NewTerminal(os.Stdin, ">>> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating terminal: %v\r\n", err)
		return
	}
	defer t.Restore()

	variables := make(map[string]float64)
	fmt.Printf("Enter an expression (or 'exit' to quit):\r\n")
	for {
		input, err := t.ReadLine()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("^c\r\n")
				return
			}
			fmt.Fprintf(os.Stderr, "error reading input: %v\r\n", err)
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "":
			continue
		case "exit":
			return
		case "history":
			fmt.Printf("%s\r\n", t.GetHistory())
			continue
		case "clear":
			t.ClearHistory()
			fmt.Printf("History cleared\r\n")
			continue
		}

		lexer, err := parser.NewLexer(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating lexer: %v\r\n", err)
			continue
		}

		expr, err := parser.NewExpressionFromLexer(lexer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating expression: %v\r\n", err)
			continue
		}

		// Handle variable assignment
		if expr.IsOPAssignment() {
			varName, rhs := expr.GetAssignment()
			if varName != "" && rhs != nil {
				val, err := rhs.Evaluate(variables)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error evaluating assignment: %v\r\n", err)
					continue
				}

				variables[varName] = val
				continue
			}
		}

		result, err := expr.Evaluate(variables)
		if err != nil {
			if err == parser.ErrNilExpression {
				continue
			}
			fmt.Fprintf(os.Stderr, "error evaluating expression: %v\r\n", err)
			continue
		}
		fmt.Printf("%g\r\n", result)
	}
}

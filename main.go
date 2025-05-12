package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"simplecalc/pkg/parser"
	"simplecalc/pkg/terminal"
)

func help() {
	msg := `Simple Calculator
Commands:
  - help: Show this help message
  - exit: Exit the calculator
  - history: Show the command history
  - clear: Clear the history
  - <expression>: Evaluate the expression
  - <var> = <expression>: Assign the expression to the variable
  - <var>: Show the value of the variable
  - <expression1>; <expression2>; ...: Evaluate multiple expressions
  - <var1> = <expression1>; <var2> = <expression2>; ...: Assign multiple variables
Examples:
  >>> 2 + 6
  >>> x = 7 + 8
  >>> y = x * 2
  >>> z = x / (2.5 * (-6 + y))
  >>> z
  >>> a = 2; b = -17; c = b / (a + -12); c
`

	// Add CRLF to each line
	lines := strings.ReplaceAll(msg, "\n", "\r\n")

	fmt.Print(lines)
}

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
		case "help":
			help()
			continue
		case "history":
			fmt.Printf("%s\r\n", t.GetHistory())
			continue
		case "clear":
			t.ClearHistory()
			fmt.Printf("History cleared\r\n")
			continue
		}

		// Support for multi-line input separated by semicolons
		lines := strings.SplitSeq(input, ";")
		for line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			lexer, err := parser.NewLexer(line)
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
}

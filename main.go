package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"simplecalc/pkg/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter an expression (or 'exit' to quit):")
	for {
		print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		lexer, err := parser.NewLexer(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		expr, err := parser.NewExpressionFromLexer(lexer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		result, err := expr.Evaluate()
		if err != nil {
			if err == parser.ErrNilExpression {
				continue
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}
		fmt.Printf("%g\n", result)
	}
}

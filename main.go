package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"simplecalc/pkg/parser"
	"simplecalc/pkg/terminal"
)

func printRusults(results []float64) {
	for _, result := range results {
		// Minimize digits and no scientific notation
		// if number is greater than 1e6-1 or less than -1e6+1
		fmt.Printf("%s\r\n", strconv.FormatFloat(result, 'f', -1, 64))
	}
}

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
  >>> a = 2; b = -17; c = -b / (a + -12); c
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

	p := parser.NewParser()

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

		// Single input may has multiple expressions separated by semicolons
		results, err := p.Parse(input)
		if err != nil {
			if err == parser.ErrNilExpression {
				continue
			}
			fmt.Fprintf(os.Stderr, "error from parser: %v\r\n", err)
			continue
		}

		printRusults(results)
	}
}

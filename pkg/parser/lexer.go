package parser

import (
	"fmt"
	"strings"
	"unicode"

	"simplecalc/pkg/parser/operator"
)

type Lexer struct {
	tokens []Token
	cursor int
}

func NewLexer(input string) (*Lexer, error) {
	if len(input) == 0 {
		return &Lexer{
			tokens: []Token{},
			cursor: 0,
		}, nil
	}

	l := Lexer{}
	err := l.parseTokens(input)
	if err != nil {
		return nil, fmt.Errorf("failed to run NewLexer: %w", err)
	}

	return &l, nil
}

func (l *Lexer) parseTokens(input string) error {
	// Remove all whitespace characters (spaces, tabs, newlines, etc.)
	input = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)

	// lexer.cursor has two purposes:
	// 1. It is used to track the current position in the input string
	// 2. It is used to track the current position in the tokens slice

	// Here we are using l.cursor to track the position in the input string
	l.cursor = 0

	for l.cursor < len(input) {
		char := input[l.cursor]

		// Return operator and new cursor if char is an operator
		op, newCursor := operator.LexWithOperator(&input, l.cursor)
		if op != nil {
			l.tokens = append(l.tokens, NewOPToken(op))
			l.cursor = newCursor
			continue
		}

		// If char is not an operator, we need to check if it's a number or a variable name
		switch char {
		case '.':
			if l.cursor < len(input)-1 && unicode.IsDigit(rune(input[l.cursor+1])) {
				err := l.readNumber(input)
				if err != nil {
					return fmt.Errorf("failed to read number: %w", err)
				}
			} else {
				return fmt.Errorf("invalid number %s", string(char))
			}
		default:
			if unicode.IsDigit(rune(char)) {
				err := l.readNumber(input)
				if err != nil {
					return fmt.Errorf("failed to read number: %w", err)
				}
			} else if char == '_' || unicode.IsLetter(rune(char)) {
				err := l.readVarName(input, false)
				if err != nil {
					return fmt.Errorf("failed to read variable name: %w", err)
				}
			} else {
				return fmt.Errorf("illegal character %s", string(char))
			}
		}
	}

	// Reset the cursor
	// Here we are using l.cursor to track the position in the tokens slice
	l.cursor = 0

	return nil
}

// Next returns the next token in the list.
// It returns EOF if there are no more tokens.
func (l *Lexer) Next() Token {
	if l.cursor >= len(l.tokens) {
		return NewEOFToken()
	}

	token := l.tokens[l.cursor]
	l.cursor++

	return token
}

func (l *Lexer) Peek() Token {
	if l.cursor >= len(l.tokens) {
		return NewEOFToken()
	}

	return l.tokens[l.cursor]
}

func (l *Lexer) HasNext() bool {
	return l.cursor < len(l.tokens)
}

func (l *Lexer) String() string {
	if len(l.tokens) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteRune('[')
	for i, token := range l.tokens {
		sb.WriteRune('"')
		sb.WriteString(token.String())
		sb.WriteRune('"')
		if i < len(l.tokens)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(']')

	return sb.String()
}

func (l *Lexer) readNumber(input string) error {
	var sb strings.Builder

	// Check for negative number
	if l.cursor < len(input) && input[l.cursor] == '-' {
		sb.WriteByte('-')
		l.cursor++
	}

	hasDecimal := false
	for l.cursor < len(input) && (unicode.IsDigit(rune(input[l.cursor])) || input[l.cursor] == '.') {
		if input[l.cursor] == '.' {
			if hasDecimal {
				// If we already have more than one decimal point, it's illegal
				// Append the number we have so far and the illegal character
				sb.WriteByte(input[l.cursor])
				return fmt.Errorf("invalid number %s", sb.String())
			}
			hasDecimal = true
		}
		sb.WriteByte(input[l.cursor])
		l.cursor++
	}

	if sb.Len() == 0 || (sb.Len() == 1 && sb.String() == "-") {
		return fmt.Errorf("invalid number: %c", input[l.cursor])
	}

	l.tokens = append(l.tokens, NewAtomNumToken(sb.String()))

	return nil
}

// readVarName reads a variable name from the input string
func (l *Lexer) readVarName(input string, negative bool) error {
	// If the variable is negative, we need to skip the '-' character
	// This is to handle cases like "-x", "(-x+3)", "-x-1".
	if negative {
		if l.cursor >= len(input)-1 {
			return fmt.Errorf("invalid variable name at end of input")
		}
		if input[l.cursor] != '-' {
			return fmt.Errorf("invalid variable name: '%c'", input[l.cursor])
		}

		l.cursor++
	}

	// Check if the variable name starts with a letter or underscore
	if l.cursor < len(input) &&
		(!unicode.IsLetter(rune(input[l.cursor])) && input[l.cursor] != '_') {
		return fmt.Errorf("invalid variable name: '%c'", input[l.cursor])
	}

	var sb strings.Builder
	if negative {
		sb.WriteByte('-')
	}
	for l.cursor < len(input) &&
		(unicode.IsLetter(rune(input[l.cursor])) ||
			unicode.IsDigit(rune(input[l.cursor])) ||
			input[l.cursor] == '_') {
		sb.WriteByte(input[l.cursor])
		l.cursor++
	}

	if sb.Len() == 0 {
		if l.cursor < len(input) {
			return fmt.Errorf("invalid variable name: %c", input[l.cursor])
		} else {
			return fmt.Errorf("invalid variable name at end of input")
		}
	}

	l.tokens = append(l.tokens, NewAtomVarToken(sb.String()))
	return nil
}

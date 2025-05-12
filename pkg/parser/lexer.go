package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	tokens []Token
	cursor int
}

func NewLexer(input string) (*Lexer, error) {
	if len(input) == 0 {
		return &Lexer{
			tokens: []Token{
				NewEOFToken(),
			},
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

	// Check if the current token is a negative number
	// This is to handle cases like "-5", "(-5+3)", "-3-1", "-.5" and "-5.3".
	isNegativeNumber := func() bool {
		// Check ether if the current character is a '-' and the next character is a digit
		if ((l.cursor < len(input)-1 && input[l.cursor] == '-' &&
			unicode.IsDigit(rune(input[l.cursor+1]))) ||

			// Or
			// If the current character is a '-' and the next character
			// is a '.' and the next character of '.' is a digit
			(l.cursor < len(input)-2 && input[l.cursor] == '-' &&
				input[l.cursor+1] == '.' &&
				unicode.IsDigit(rune(input[l.cursor+2])))) &&

			// And must be
			// Check the current char is either the start of the input
			// or the previous token is a left parenthesis or an arithmetic operator
			// or an assignment operator
			(l.cursor == 0 || l.tokens[len(l.tokens)-1].IsOPLeftParen() ||
				l.tokens[len(l.tokens)-1].IsArithmeticOperator() ||
				l.tokens[len(l.tokens)-1].IsOPAssign()) {

			return true
		}

		return false
	}

	// lexer.cursor has two purposes:
	// 1. It is used to track the current position in the input string
	// 2. It is used to track the current position in the tokens slice

	// Here we are using l.cursor to track the position in the input string
	l.cursor = 0

	for l.cursor < len(input) {
		char := input[l.cursor]
		switch char {
		case '=':
			l.tokens = append(l.tokens, NewOPToken(TokenAssign, string(char)))
			l.cursor++
		case '+':
			l.tokens = append(l.tokens, NewOPToken(TokenPlus, string(char)))
			l.cursor++
		case '-':
			// Check if a negative number
			if isNegativeNumber() {
				err := l.readNumber(input)
				if err != nil {
					return fmt.Errorf("failed to read number: %w", err)
				}
				continue
			}
			// If not a negative number, treat it as a minus operator
			l.tokens = append(l.tokens, NewOPToken(TokenMinus, string(char)))
			l.cursor++
		case '*':
			l.tokens = append(l.tokens, NewOPToken(TokenMultiply, string(char)))
			l.cursor++
		case '/':
			l.tokens = append(l.tokens, NewOPToken(TokenDivide, string(char)))
			l.cursor++
		case '(':
			l.tokens = append(l.tokens, NewOPToken(TokenLeftParen, string(char)))
			l.cursor++
		case ')':
			l.tokens = append(l.tokens, NewOPToken(TokenRightParen, string(char)))
			l.cursor++
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
				err := l.readVarName(input)
				if err != nil {
					return fmt.Errorf("failed to read variable name: %w", err)
				}
			} else {
				return fmt.Errorf("illegal character %s", string(char))
			}
		}
	}

	l.tokens = append(l.tokens, NewEOFToken())

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
func (l *Lexer) readVarName(input string) error {
	// Check if the variable name starts with a letter or underscore
	if l.cursor < len(input) &&
		(!unicode.IsLetter(rune(input[l.cursor])) && input[l.cursor] != '_') {
		return fmt.Errorf("invalid variable name: '%c'", input[l.cursor])
	}

	var sb strings.Builder
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

package parser

import (
	"fmt"
	"strconv"
)

type TokenType uint8

const (
	TokenEOF TokenType = iota
	TokenAtom
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLeftParen
	TokenRightParen
	TokenAssign
)

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenAtom:
		return "Atom"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenAssign:
		return "="
	default:
		return "Unknown"
	}
}

var (
	operators = map[TokenType]bool{
		TokenPlus:       true,
		TokenMinus:      true,
		TokenMultiply:   true,
		TokenDivide:     true,
		TokenLeftParen:  true,
		TokenRightParen: true,
	}

	arithmeticOperators = map[TokenType]bool{
		TokenPlus:     true,
		TokenMinus:    true,
		TokenMultiply: true,
		TokenDivide:   true,
	}
)

type Token struct {
	typ        TokenType
	literal    string
	isVariable bool
}

// NewAtomVarToken creates a new variable atom token.
func NewAtomVarToken(literal string) Token {
	return Token{
		typ:        TokenAtom,
		literal:    literal,
		isVariable: true,
	}
}

// NewAtomNumToken creates a new number atom token.
func NewAtomNumToken(literal string) Token {
	return Token{
		typ:     TokenAtom,
		literal: literal,
	}
}

// NewOPToken creates a new operator token.
func NewOPToken(tokenType TokenType, literal string) Token {
	return Token{
		typ:     tokenType,
		literal: literal,
	}
}

// NewEOFToken creates a new EOF token.
func NewEOFToken() Token {
	return Token{
		typ: TokenEOF,
	}
}

func (t Token) IsAtom() bool {
	return t.typ == TokenAtom
}

func (t Token) IsOperator() bool {
	return operators[t.typ]
}

func (t Token) IsAtomVariable() bool {
	return t.isVariable
}

func (t Token) IsArithmeticOperator() bool {
	return arithmeticOperators[t.typ]
}

func (t Token) IsOPLeftParen() bool {
	return t.typ == TokenLeftParen
}

func (t Token) IsOPRightParen() bool {
	return t.typ == TokenRightParen
}

func (t Token) IsOPAssign() bool {
	return t.typ == TokenAssign
}

func (t Token) IsEOF() bool {
	return t.typ == TokenEOF
}

// GetVarName returns the variable name of the token
func (t Token) GetVarName() string {
	if !t.isVariable {
		return ""
	}

	if t.literal == "" {
		// Must be a bug from the lexer, don't recover it
		panic("variable name is empty")
	}

	return t.literal
}

func (t Token) GetValue() float64 {
	if t.typ != TokenAtom || t.isVariable {
		return 0
	}

	value, err := strconv.ParseFloat(t.literal, 64)
	if err != nil {
		// Must be a bug, don't recover it
		// This should not happen as we are already checking the type
		// and the literal should be a valid number.
		panic(fmt.Sprintf("failed to parse token value: %s", err))
	}

	return value
}

func (t Token) GetType() TokenType {
	return t.typ
}

func (t Token) IsInfixOperator() bool {
	switch t.typ {
	case TokenAssign, TokenPlus, TokenMinus, TokenMultiply, TokenDivide:
		return true
	default:
		return false
	}
}

func (t Token) IsPrefixOperator() bool {
	switch t.typ {
	case TokenPlus, TokenMinus:
		return true
	default:
		return false
	}
}

// GetInfixBindingPower returns the left and right binding powers of an infix operator.
// Two different binding powers determine the ordering behavior to ensure predictable
// and testable parsing.
func (t Token) GetInfixBindingPower() (float32, float32, error) {
	switch t.typ {
	case TokenAssign:
		return 0.2, 0.1, nil
	case TokenPlus, TokenMinus:
		return 1.0, 1.1, nil
	case TokenMultiply, TokenDivide:
		return 2.0, 2.1, nil
	default:
		return 0, 0, fmt.Errorf("unknown infix operator: '%s'", t.typ)
	}
}

// GetPrefixBindingPower returns the right-hand side binding power
// of the prefix operator, because it's right associative only
func (t Token) GetPrefixBindingPower() (float32, error) {
	switch t.typ {
	case TokenPlus, TokenMinus:
		return 3.0, nil
	default:
		return 0, fmt.Errorf("unknown prefix operator: '%s'", t.typ)
	}
}

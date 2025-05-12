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
	typ     TokenType
	literal string
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

func (t Token) IsArithmeticOperator() bool {
	return arithmeticOperators[t.typ]
}

func (t Token) IsOPLeftParen() bool {
	return t.typ == TokenLeftParen
}

func (t Token) IsOPRightParen() bool {
	return t.typ == TokenRightParen
}

func (t Token) IsEOF() bool {
	return t.typ == TokenEOF
}

func (t Token) GetValue() float64 {
	if t.typ != TokenAtom {
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

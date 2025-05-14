package parser

import (
	"fmt"
	"strconv"

	"simplecalc/pkg/parser/operator"
)

type TokenType uint8

const (
	TokenEOF TokenType = iota
	TokenAtom
	TokenOperator
)

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenAtom:
		return "Atom"
	case TokenOperator:
		return "Operator"
	default:
		return "Unknown"
	}
}

type Token struct {
	typ        TokenType
	literal    string
	operator   operator.Operator
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

func (t Token) String() string {
	if t.IsEOF() {
		return TokenEOF.String()
	}

	return t.literal
}

// NewAtomNumToken creates a new number atom token.
func NewAtomNumToken(literal string) Token {
	if literal == "" {
		panic("literal is empty")
	}

	return Token{
		typ:     TokenAtom,
		literal: literal,
	}
}

// NewOPToken creates a new operator token.
func NewOPToken(op operator.Operator) Token {
	return Token{
		typ:      TokenOperator,
		operator: op,
		literal:  op.GetLiteral(),
	}
}

// NewOPTokenByLiteral creates a new operator token by literal.
func NewOPTokenByLiteral(literal string) Token {
	return Token{
		typ:      TokenOperator,
		operator: operator.GetOperator(literal),
		literal:  literal,
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
	return t.operator != nil
}

func (t Token) IsAtomVariable() bool {
	return t.isVariable
}

func (t Token) IsArithmeticOperator() bool {
	if t.typ != TokenOperator {
		panic(fmt.Sprintf("'%s' (%s) is not a operator", t.literal, t.typ))
	}

	return t.operator.IsArithmeticOperator()
}

func (t Token) IsTheOperator(literal string) bool {
	if t.typ != TokenOperator {
		panic(fmt.Sprintf("'%s' (%s) is not a operator", t.literal, t.typ))
	}

	return t.operator.Is(literal)
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

func (t Token) GetOperator() operator.Operator {
	if t.typ != TokenOperator {
		panic(fmt.Sprintf("'%s' (%s) is not a operator", t.literal, t.typ))
	}

	return t.operator
}

func (t Token) IsInfixOperator() bool {
	if t.typ != TokenOperator {
		panic(fmt.Sprintf("'%s' (%s) is not a operator", t.literal, t.typ))
	}

	return t.operator.IsInfixOperator()
}

func (t Token) IsPrefixOperator() bool {
	if t.typ != TokenOperator {
		panic(fmt.Sprintf("'%s' (%s) is not a operator", t.literal, t.typ))
	}

	return t.operator.IsPrefixOperator()
}

// GetInfixBindingPower returns the left and right binding powers of an infix operator.
// Two different binding powers determine the ordering behavior to ensure predictable
// and testable parsing.
func (t Token) GetInfixBindingPower() (float32, float32, error) {
	if t.typ != TokenOperator {
		return 0, 0, fmt.Errorf("'%s' is not a operator", t.literal)
	}

	return t.operator.GetInfixBindingPower()
}

// GetPrefixBindingPower returns the right-hand side binding power
// of the prefix operator, because it's right associative only
func (t Token) GetPrefixBindingPower() (float32, error) {
	if t.typ != TokenOperator {
		return 0, fmt.Errorf("'%s' is not a operator", t.literal)
	}

	return t.operator.GetPrefixBindingPower()
}

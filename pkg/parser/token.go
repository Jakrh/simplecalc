package parser

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
	operations = map[TokenType]bool{
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
	Type    TokenType
	Literal string
}

func (t Token) IsAtom() bool {
	return t.Type == TokenAtom
}

func (t Token) IsOperator() bool {
	return operations[t.Type]
}

func (t Token) IsArithmeticOperator() bool {
	return arithmeticOperators[t.Type]
}

func (t Token) IsLeftParen() bool {
	return t.Type == TokenLeftParen
}

func (t Token) IsRightParen() bool {
	return t.Type == TokenRightParen
}

func (t Token) IsEOF() bool {
	return t.Type == TokenEOF
}

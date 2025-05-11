package parser

import (
	"fmt"
	"strconv"
)

type ExprType uint8

const (
	ExprTypeAtomic ExprType = iota
	ExprTypeOperation
)

var (
	ErrNilExpression  = fmt.Errorf("expression is nil")
	ErrDivisionByZero = fmt.Errorf("division by zero")
)

type Expression struct {
	Type  ExprType
	Value float64
	OP    TokenType
	Left  *Expression
	Right *Expression
}

func (e *Expression) GetType() ExprType {
	return e.Type
}

func (e *Expression) Evaluate() (float64, error) {
	if e == nil {
		return 0, ErrNilExpression
	}

	if e.Type == ExprTypeAtomic {
		return e.Value, nil
	}

	var leftValue, rightValue float64
	var err error
	if e.Left != nil {
		leftValue, err = e.Left.Evaluate()
		if err != nil {
			return 0, fmt.Errorf("failed to evaluate left expression: %w", err)
		}
	}
	if e.Right != nil {
		rightValue, err = e.Right.Evaluate()
		if err != nil {
			return 0, fmt.Errorf("failed to evaluate right expression: %w", err)
		}
	}
	switch e.OP {
	case TokenPlus:
		return leftValue + rightValue, nil
	case TokenMinus:
		return leftValue - rightValue, nil
	case TokenMultiply:
		return leftValue * rightValue, nil
	case TokenDivide:
		if rightValue == 0 {
			return 0, ErrDivisionByZero
		}
		return leftValue / rightValue, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", e.OP)
	}
}

func (e *Expression) String() string {
	if e == nil {
		return ""
	}

	if e.Type == ExprTypeAtomic {
		return strconv.FormatFloat(e.Value, 'f', -1, 64)
	} else {
		return fmt.Sprintf("(%s %s %s)", e.OP, e.Left, e.Right)
	}
}

func NewExpressionFromLexer(lexer *Lexer) (*Expression, error) {
	expr, err := parseExpressions(lexer, 0.0)
	if err != nil {
		return nil, fmt.Errorf("failed to run NewExpressionFromLexer: %w", err)
	}

	return expr, nil
}

func newAtomicExpression(value float64) *Expression {
	return &Expression{
		Type:  ExprTypeAtomic,
		Value: value,
	}
}

func newOperationExpression(op TokenType, left, right *Expression) *Expression {
	return &Expression{
		Type:  ExprTypeOperation,
		OP:    op,
		Left:  left,
		Right: right,
	}
}

func parseExpressions(lexer *Lexer, minBP float32) (*Expression, error) {
	var lhs *Expression
	lhsToken := lexer.Next()
	if lhsToken.IsAtom() {
		value, err := strconv.ParseFloat(lhsToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %w", err)
		}
		lhs = newAtomicExpression(value)
	} else if lhsToken.IsOperator() {
		if lhsToken.Type == TokenLeftParen {
			var err error
			lhs, err = parseExpressions(lexer, 0.0)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression: %w", err)
			}
			if lexer.Next().Type != TokenRightParen {
				return nil, fmt.Errorf("missing right parenthesis")
			}
		}
	} else if lhsToken.IsEOF() {
		return nil, nil
	} else {
		return nil, fmt.Errorf("unexpected token: %s", lhsToken.Literal)
	}

	for {
		op := lexer.Peek()
		if op.Type == TokenEOF || op.Type == TokenRightParen {
			break
		}
		lBP, rBP := infixBindingPower(op.Type)
		if lBP < minBP {
			break
		}

		lexer.Next() // consume the operator token
		rhs, err := parseExpressions(lexer, rBP)
		if err != nil {
			return nil, fmt.Errorf("failed to parse right-hand side: %w", err)
		}
		lhs = newOperationExpression(op.Type, lhs, rhs)
	}

	return lhs, nil
}

func infixBindingPower(op TokenType) (float32, float32) {
	switch op {
	case TokenAssign:
		return 0.2, 0.1
	case TokenPlus, TokenMinus:
		return 1.0, 1.1
	case TokenMultiply, TokenDivide:
		return 2.0, 2.1
	default:
		panic(fmt.Sprintf("unknown operator: '%s'", op))
	}
}

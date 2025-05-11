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
	ErrNilExpression           = fmt.Errorf("expression is nil")
	ErrDivisionByZero          = fmt.Errorf("division by zero")
	ErrMissingLeftParenthesis  = fmt.Errorf("missing left parenthesis")
	ErrMissingRightParenthesis = fmt.Errorf("missing right parenthesis")
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
	// parenBalance is used to track the balance of parentheses.
	// Increment it when we encounter a left parenthesis
	// and decrement it when we encounter a right parenthesis.
	// If it goes below zero, we have some unmatched right parentheses.
	parenBalance := 0

	var parse func(*Lexer, float32) (*Expression, error)
	parse = func(lexer *Lexer, minBP float32) (*Expression, error) {
		var lhs *Expression
		lhsToken := lexer.Next()
		if lhsToken.IsAtom() {
			value, err := strconv.ParseFloat(lhsToken.Literal, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse an atom: %w", err)
			}
			lhs = newAtomicExpression(value)
		} else if lhsToken.IsOperator() {
			// Handle parentheses and EOF tokens
			if lhsToken.IsLeftParen() {
				// If we encounter a left parenthesis, we need to parse
				// the inside expression recursively as left-hand side
				// and check for the right parenthesis.
				parenBalance++
				var err error
				lhs, err = parse(lexer, 0.0)
				if err != nil {
					return nil, fmt.Errorf("failed to parse expression: %w", err)
				}
				if lexer.Next().IsRightParen() {
					parenBalance--
				} else if parenBalance > 0 {
					// Return an error if we don't find a matching right parenthesis
					return nil, ErrMissingRightParenthesis
				}
			} else if lhsToken.IsRightParen() {
				parenBalance--
			}
		} else if lhsToken.IsEOF() {
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected token: %s", lhsToken.Literal)
		}

		for {
			// Handle EOF and right parenthesis tokens
			op := lexer.Peek()
			if op.IsEOF() {
				break
			} else if op.IsRightParen() {
				// Return an error if we find a right parenthesis
				// without a matching left parenthesis
				if parenBalance == 0 {
					return nil, ErrMissingLeftParenthesis
				}
				break
			}

			// Stop parsing the right-hand side expression if the left-hand side
			// binding power of this operator is less than the minimum binding power
			lBP, rBP := infixBindingPower(op.Type)
			if lBP < minBP {
				break
			}

			// Parse the right-hand side
			lexer.Next() // Consume the operator token
			rhs, err := parse(lexer, rBP)
			if err != nil {
				return nil, fmt.Errorf("failed to parse right-hand side: %w", err)
			}
			lhs = newOperationExpression(op.Type, lhs, rhs)
		}

		return lhs, nil
	}

	expr, err := parse(lexer, minBP)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}

	// Return an error if we have unmatched right parentheses
	if parenBalance < 0 {
		return nil, ErrMissingLeftParenthesis
	}

	return expr, nil
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

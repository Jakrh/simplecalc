package parser

import (
	"fmt"
	"math"
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
	ErrNumOutOfRange           = fmt.Errorf("number is too large/small that lost percision in float64")
)

type Expression struct {
	typ          ExprType
	value        float64
	variableName string
	op           TokenType
	left         *Expression
	right        *Expression
}

func (e *Expression) GetType() ExprType {
	return e.typ
}

func (e *Expression) Evaluate(variables map[string]float64) (float64, error) {
	value, err := e.evaluate(variables)
	if err != nil {
		return 0, err
	}

	// Check if the value is too large/small
	if e.isNumOutOfRange(value) {
		return 0, ErrNumOutOfRange
	}

	return value, nil
}

func (e *Expression) evaluate(variables map[string]float64) (float64, error) {
	if e == nil {
		return 0, ErrNilExpression
	}

	// If the expression is atomic, get value
	// from the value field or from the variables map
	if e.IsAtom() {
		if e.IsAtomVarName() {
			varName := e.GetVarName()

			// Check if the variable name starts with a negative sign
			negative := false
			if len(varName) > 0 && varName[0] == '-' {
				varName = varName[1:] // Remove the negative sign
				negative = true       // Set the negative flag
			}

			if val, ok := variables[varName]; ok {
				if negative {
					val = -val
				}

				return val, nil
			}
			return 0, fmt.Errorf("undefined variable '%s'", varName)
		}

		return e.value, nil
	}

	// If the expression is an operation, evaluate the left and right expressions
	var leftValue, rightValue float64
	var err error
	if e.left != nil {
		leftValue, err = e.left.Evaluate(variables)
		if err != nil {
			return 0, fmt.Errorf("failed to evaluate left expression: %w", err)
		}
	}
	if e.right != nil {
		rightValue, err = e.right.Evaluate(variables)
		if err != nil {
			return 0, fmt.Errorf("failed to evaluate right expression: %w", err)
		}
	}
	switch e.op {
	case TokenAssign:
		return 0, nil
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
	case TokenPower:
		return math.Pow(leftValue, rightValue), nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", e.op)
	}
}

func (e *Expression) String() string {
	if e == nil {
		return ""
	}

	if e.IsAtom() {
		if e.IsAtomVarName() {
			return e.GetVarName()
		} else {
			return strconv.FormatFloat(e.value, 'f', -1, 64)
		}
	} else {
		return fmt.Sprintf("(%s %s %s)", e.op, e.left, e.right)
	}
}

func NewExpressionFromLexer(lexer *Lexer) (*Expression, error) {
	expr, err := parseExpressions(lexer, 0.0)
	if err != nil {
		return nil, fmt.Errorf("failed to run NewExpressionFromLexer: %w", err)
	}

	return expr, nil
}

func newAtomicVarExpression(literal string) *Expression {
	return &Expression{
		typ:          ExprTypeAtomic,
		variableName: literal,
	}
}

func newAtomicNumExpression(value float64) *Expression {
	return &Expression{
		typ:   ExprTypeAtomic,
		value: value,
	}
}

func newOperationExpression(op TokenType, left, right *Expression) *Expression {
	return &Expression{
		typ:   ExprTypeOperation,
		op:    op,
		left:  left,
		right: right,
	}
}

func (e *Expression) IsAtom() bool {
	return e != nil && e.typ == ExprTypeAtomic
}

func (e *Expression) IsOperation() bool {
	return e != nil && e.typ == ExprTypeOperation
}

// IsAtomVarName checks if the expression is an atom variable name.
func (e *Expression) IsAtomVarName() bool {
	return e != nil && e.IsAtom() && e.variableName != ""
}

func (e *Expression) IsOPAssignment() bool {
	return e != nil && e.IsOperation() && e.op == TokenAssign
}

// GetVarName returns the variable name of the expression
func (e *Expression) GetVarName() string {
	if e.IsAtomVarName() {
		return e.variableName
	}

	return ""
}

// GetAssignment returns the right-hand side expression
// of the operation from an assignment expression.
func (e *Expression) GetAssignment() (string, *Expression) {
	if e.IsOperation() && e.IsOPAssignment() {
		if e.left == nil || !e.left.IsAtom() {
			// Must be a bug, don't recover it
			panic("left expression is not an atom")
		}

		varName := e.left.variableName
		if varName == "" {
			// Must be a bug from the lexer, don't recover it
			// This atom may be a number not a variable name
			panic("variable name is empty")
		}

		expr := e.right
		if expr == nil {
			// Must be a bug, don't recover it
			panic("right expression is nil")
		}

		return varName, expr
	}

	return "", nil
}

func (e *Expression) isNumOutOfRange(value float64) bool {
	const effectiveBoundary = float64(1 << 53)
	if value <= -effectiveBoundary || value >= effectiveBoundary {
		return true
	}

	return false
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
			if lhsToken.IsAtomVariable() {
				varName := lhsToken.GetVarName()
				if varName == "" {
					return nil, fmt.Errorf("variable name is empty")
				}
				lhs = newAtomicVarExpression(varName)
			} else {
				lhs = newAtomicNumExpression(lhsToken.GetValue())
			}
		} else if lhsToken.IsOperator() {
			// Handle parentheses and EOF tokens
			if lhsToken.IsOPLeftParen() {
				// If we encounter a left parenthesis, we need to parse
				// the inside expression recursively as left-hand side
				// and check for the right parenthesis.
				parenBalance++
				var err error
				lhs, err = parse(lexer, 0.0)
				if err != nil {
					return nil, fmt.Errorf("failed to parse expression: %w", err)
				}
				if lexer.Next().IsOPRightParen() {
					parenBalance--
				} else if parenBalance > 0 {
					// Return an error if we don't find a matching right parenthesis
					return nil, ErrMissingRightParenthesis
				}
			} else if lhsToken.IsOPRightParen() {
				parenBalance--
			} else if lhsToken.IsPrefixOperator() {
				rBP, err := lhsToken.GetPrefixBindingPower()
				if err != nil {
					return nil, fmt.Errorf("failed to get prefix binding power: %w", err)
				}
				lhs, err = parse(lexer, rBP)
				if err != nil {
					return nil, fmt.Errorf("failed to parse expression: %w", err)
				}
				if lhs == nil {
					return nil,
						fmt.Errorf("missing left-hand side expression from prefix operator: %s",
							lhsToken.GetType())
				}
				// Add a new operation expression with the prefix operator and 0 as the left operand,
				// the parsed expression as the right operand.
				// This is to handle cases like "-x" as "(- 0 x) and "+x" as "(+ 0 x)".
				lhs = newOperationExpression(lhsToken.GetType(), newAtomicNumExpression(0), lhs)
			}
		} else if lhsToken.IsEOF() {
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected token: %s", lhsToken.literal)
		}

		for {
			// Handle EOF and right parenthesis tokens
			op := lexer.Peek()
			if op.IsEOF() {
				break
			} else if op.IsOPRightParen() {
				// Return an error if we find a right parenthesis
				// without a matching left parenthesis
				if parenBalance == 0 {
					return nil, ErrMissingLeftParenthesis
				}
				break
			}

			// Stop parsing the right-hand side expression if the left-hand side
			// binding power of this operator is less than the minimum binding power
			lBP, rBP, err := op.GetInfixBindingPower()
			if err != nil {
				return nil, fmt.Errorf("failed to get binding power: %w", err)
			}
			if lBP < minBP {
				break
			}

			// Parse the right-hand side
			lexer.Next() // Consume the operator token
			rhs, err := parse(lexer, rBP)
			if err != nil {
				return nil, fmt.Errorf("failed to parse right-hand side: %w", err)
			}
			lhs = newOperationExpression(op.GetType(), lhs, rhs)
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

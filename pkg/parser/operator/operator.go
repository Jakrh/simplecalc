package operator

import (
	"fmt"
)

var (
	ErrDivisionByZero      = fmt.Errorf("division by zero")
	ErrInvalidOperator     = fmt.Errorf("invalid operator")
	ErrInvalidOperand      = fmt.Errorf("invalid operand")
	ErrOperatorNotFound    = fmt.Errorf("operator not found")
	ErrInvalidOperandCount = fmt.Errorf("invalid operand count")
	ErrNotInfixOperator    = fmt.Errorf("not infix operator")
	ErrNotPrefixOperator   = fmt.Errorf("not prefix operator")
)

type OPid uint8

type Operator interface {
	Is(literal string) bool
	IsArithmeticOperator() bool
	IsInfixOperator() bool
	IsPrefixOperator() bool
	isGroupingOperator() bool
	GetLiteral() string
	GetInfixBindingPower() (float32, float32, error)
	GetPrefixBindingPower() (float32, error)
	Lex(input *string, cursor int) (string, int)
	Evaluate(oprands []float64) (float64, error)
	String() string
}

// set containers that registers all operators from init
var (
	allOperators = map[string]Operator{}
	lexOperators = map[byte]map[int]Operator{}
)

func registerOperator(op Operator) {
	if _, ok := allOperators[op.GetLiteral()]; ok {
		panic(fmt.Sprintf("operator '%s' already registered", op.GetLiteral()))
	}
	allOperators[op.GetLiteral()] = op

	opLiteral := op.GetLiteral()
	if len(opLiteral) == 0 {
		panic("operator literal cannot be empty")
	}
	opStartWith := opLiteral[0]
	if _, ok := lexOperators[opStartWith]; !ok {
		lexOperators[opStartWith] = make(map[int]Operator)
	}
	lexOperators[opStartWith][len(opLiteral)] = op
}

// LexWithOperator lexes the input string and returns the operator and the new cursor position
// It expects the input string to not contain any whitespace characters!
// It won't check and panic with whitespace for performance
func LexWithOperator(input *string, cursor int) (Operator, int) {
	if cursor >= len(*input) {
		return nil, cursor
	}

	char := (*input)[cursor]
	if ops, ok := lexOperators[char]; ok {
		for i := len(ops); i > 0; i-- {
			if op, ok := ops[i]; ok {
				token, newCursor := op.Lex(input, cursor)
				if token != "" {
					return op, newCursor
				}
				continue
			}
		}
		// No operator found
		return nil, cursor
	} else {
		// No operator found
		cursor++
		return nil, cursor
	}
}

func GetOperator(literal string) Operator {
	if op, ok := allOperators[literal]; ok {
		return op
	}

	panic(fmt.Errorf("%w: '%s'", ErrOperatorNotFound, literal))
}

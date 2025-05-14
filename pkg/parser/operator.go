package parser

import (
	"fmt"
	"math"
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

// --------------------------------------------------------------

type add struct {
	literal string
}

func init() {
	registerOperator(&add{
		literal: "+",
	})
}

func (o *add) Is(literal string) bool {
	return o.literal == literal
}

func (o *add) IsArithmeticOperator() bool {
	return true
}

func (o *add) IsInfixOperator() bool {
	return true
}

func (o *add) IsPrefixOperator() bool {
	return true
}

func (o *add) isGroupingOperator() bool {
	return false
}

func (o *add) GetLiteral() string {
	return o.literal
}

func (o *add) GetInfixBindingPower() (float32, float32, error) {
	return 1.0, 1.1, nil
}

func (o *add) GetPrefixBindingPower() (float32, error) {
	return 3.0, nil
}

func (o *add) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

func (o *add) Evaluate(oprands []float64) (float64, error) {
	if len(oprands) != 2 {
		return 0,
			fmt.Errorf(
				"%w: must have exactly 2 operands for '%s' operator",
				ErrInvalidOperandCount,
				o.literal)
	}

	return oprands[0] + oprands[1], nil
}

func (o *add) String() string {
	return o.literal
}

// --------------------------------------------------------------

type minus struct {
	literal string
}

func init() {
	registerOperator(&minus{
		literal: "-",
	})
}

func (o *minus) Is(literal string) bool {
	return o.literal == literal
}

func (o *minus) IsArithmeticOperator() bool {
	return true
}

func (o *minus) IsInfixOperator() bool {
	return true
}

func (o *minus) IsPrefixOperator() bool {
	return true
}

func (o *minus) isGroupingOperator() bool {
	return false
}

func (o *minus) GetLiteral() string {
	return o.literal
}

func (o *minus) GetInfixBindingPower() (float32, float32, error) {
	return 1.0, 1.1, nil
}

func (o *minus) GetPrefixBindingPower() (float32, error) {
	return 3.0, nil
}

func (o *minus) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

func (o *minus) Evaluate(oprands []float64) (float64, error) {
	if len(oprands) != 2 {
		return 0,
			fmt.Errorf(
				"%w: must have exactly 2 operands for '%s' operator",
				ErrInvalidOperandCount,
				o.literal)
	}

	return oprands[0] - oprands[1], nil
}

func (o *minus) String() string {
	return o.literal
}

// --------------------------------------------------------------

type multiply struct {
	literal string
}

func init() {
	registerOperator(&multiply{
		literal: "*",
	})
}

func (o *multiply) Is(literal string) bool {
	return o.literal == literal
}

func (o *multiply) IsArithmeticOperator() bool {
	return true
}

func (o *multiply) IsInfixOperator() bool {
	return true
}

func (o *multiply) IsPrefixOperator() bool {
	return false
}

func (o *multiply) isGroupingOperator() bool {
	return false
}

func (o *multiply) GetLiteral() string {
	return o.literal
}

func (o *multiply) GetInfixBindingPower() (float32, float32, error) {
	return 2.0, 2.1, nil
}

func (o *multiply) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *multiply) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

func (o *multiply) Evaluate(oprands []float64) (float64, error) {
	if len(oprands) != 2 {
		return 0,
			fmt.Errorf(
				"%w: must have exactly 2 operands for '%s' operator",
				ErrInvalidOperandCount,
				o.literal)
	}

	return oprands[0] * oprands[1], nil
}

func (o *multiply) String() string {
	return o.literal
}

// --------------------------------------------------------------

type divide struct {
	literal string
}

func init() {
	registerOperator(&divide{
		literal: "/",
	})
}

func (o *divide) Is(literal string) bool {
	return o.literal == literal
}

func (o *divide) IsArithmeticOperator() bool {
	return true
}

func (o *divide) IsInfixOperator() bool {
	return true
}

func (o *divide) IsPrefixOperator() bool {
	return false
}

func (o *divide) isGroupingOperator() bool {
	return false
}

func (o *divide) GetLiteral() string {
	return o.literal
}

func (o *divide) GetInfixBindingPower() (float32, float32, error) {
	return 2.0, 2.1, nil
}

func (o *divide) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *divide) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

func (o *divide) Evaluate(oprands []float64) (float64, error) {
	if len(oprands) != 2 {
		return 0,
			fmt.Errorf(
				"%w: must have exactly 2 operands for '%s' operator",
				ErrInvalidOperandCount,
				o.literal)
	}

	if oprands[1] == 0 {
		return 0, ErrDivisionByZero
	}

	return oprands[0] / oprands[1], nil
}

func (o *divide) String() string {
	return o.literal
}

// --------------------------------------------------------------

type power struct {
	literal string
}

func init() {
	registerOperator(&power{
		literal: "**",
	})
}

func (o *power) Is(literal string) bool {
	return o.literal == literal
}

func (o *power) IsArithmeticOperator() bool {
	return true
}

func (o *power) IsInfixOperator() bool {
	return true
}

func (o *power) IsPrefixOperator() bool {
	return false
}

func (o *power) isGroupingOperator() bool {
	return false
}

func (o *power) GetLiteral() string {
	return o.literal
}

func (o *power) GetInfixBindingPower() (float32, float32, error) {
	return 4.0, 4.1, nil
}

func (o *power) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *power) Lex(input *string, cursor int) (string, int) {
	if cursor < len(*input)-1 && (*input)[cursor+1] == '*' {
		return o.literal, cursor + 2
	}

	return "", cursor
}

func (o *power) Evaluate(oprands []float64) (float64, error) {
	if len(oprands) != 2 {
		return 0,
			fmt.Errorf(
				"%w: must have exactly 2 operands for '%s' operator",
				ErrInvalidOperandCount,
				o.literal)
	}

	return math.Pow(oprands[0], oprands[1]), nil
}

func (o *power) String() string {
	return o.literal
}

// --------------------------------------------------------------

type leftParen struct {
	literal string
}

func init() {
	registerOperator(&leftParen{
		literal: "(",
	})
}

func (o *leftParen) Is(literal string) bool {
	return o.literal == literal
}

func (o *leftParen) IsArithmeticOperator() bool {
	return false
}

func (o *leftParen) IsInfixOperator() bool {
	return false
}

func (o *leftParen) IsPrefixOperator() bool {
	return false
}

func (o *leftParen) isGroupingOperator() bool {
	return true
}

func (o *leftParen) GetLiteral() string {
	return o.literal
}

func (o *leftParen) GetInfixBindingPower() (float32, float32, error) {
	return 0, 0, fmt.Errorf("%w: '%s'", ErrNotInfixOperator, o.literal)
}

func (o *leftParen) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *leftParen) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

// Evaluate is not applicable for LeftParen operator
func (o *leftParen) Evaluate(oprands []float64) (float64, error) {
	return 0, nil
}

func (o *leftParen) String() string {
	return o.literal
}

// --------------------------------------------------------------

type rightParen struct {
	literal string
}

func init() {
	registerOperator(&rightParen{
		literal: ")",
	})
}

func (o *rightParen) Is(literal string) bool {
	return o.literal == literal
}

func (o *rightParen) IsArithmeticOperator() bool {
	return false
}

func (o *rightParen) IsInfixOperator() bool {
	return false
}

func (o *rightParen) IsPrefixOperator() bool {
	return false
}

func (o *rightParen) isGroupingOperator() bool {
	return true
}

func (o *rightParen) GetLiteral() string {
	return o.literal
}

func (o *rightParen) GetInfixBindingPower() (float32, float32, error) {
	return 0, 0, fmt.Errorf("%w: '%s'", ErrNotInfixOperator, o.literal)
}

func (o *rightParen) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *rightParen) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

// Evaluate is not applicable for RightParen operator
func (o *rightParen) Evaluate(oprands []float64) (float64, error) {
	return 0, nil
}

func (o *rightParen) String() string {
	return o.literal
}

// --------------------------------------------------------------

type assign struct {
	literal string
}

func init() {
	registerOperator(&assign{
		literal: "=",
	})
}

func (o *assign) Is(literal string) bool {
	return o.literal == literal
}

func (o *assign) IsArithmeticOperator() bool {
	return false
}

func (o *assign) IsInfixOperator() bool {
	return true
}

func (o *assign) IsPrefixOperator() bool {
	return false
}

func (o *assign) isGroupingOperator() bool {
	return false
}

func (o *assign) GetLiteral() string {
	return o.literal
}

func (o *assign) GetInfixBindingPower() (float32, float32, error) {
	return 0.2, 0.1, nil
}

func (o *assign) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *assign) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

// Evaluate is not applicable for Assign operator
func (o *assign) Evaluate(oprands []float64) (float64, error) {
	return 0, nil
}

func (o *assign) String() string {
	return o.literal
}

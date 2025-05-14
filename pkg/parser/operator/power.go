package operator

import (
	"fmt"
	"math"
)

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

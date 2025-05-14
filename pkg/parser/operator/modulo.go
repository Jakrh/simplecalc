package operator

import (
	"fmt"
	"math"
)

type modulo struct {
	literal string
}

func init() {
	registerOperator(&modulo{
		literal: "%",
	})
}

func (o *modulo) Is(literal string) bool {
	return o.literal == literal
}

func (o *modulo) IsArithmeticOperator() bool {
	return true
}

func (o *modulo) IsInfixOperator() bool {
	return true
}

func (o *modulo) IsPrefixOperator() bool {
	return false
}

func (o *modulo) isGroupingOperator() bool {
	return false
}

func (o *modulo) GetLiteral() string {
	return o.literal
}

func (o *modulo) GetInfixBindingPower() (float32, float32, error) {
	return 2.0, 2.1, nil
}

func (o *modulo) GetPrefixBindingPower() (float32, error) {
	return 0, fmt.Errorf("%w: '%s'", ErrNotPrefixOperator, o.literal)
}

func (o *modulo) Lex(input *string, cursor int) (string, int) {
	return o.literal, cursor + 1
}

func (o *modulo) Evaluate(oprands []float64) (float64, error) {
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

	rounded0 := math.Round(oprands[0])
	if math.Abs(oprands[0]-rounded0) > 1e-9 {
		return 0,
			fmt.Errorf(
				"%w: modulo operator requires integer operands, got %g",
				ErrInvalidOperand,
				oprands[0])
	}
	rounded1 := math.Round(oprands[1])
	if math.Abs(oprands[1]-rounded1) > 1e-9 {
		return 0,
			fmt.Errorf(
				"%w: modulo operator requires integer operands, got %g",
				ErrInvalidOperand,
				oprands[1])
	}

	return float64(int(rounded0) % int(rounded1)), nil
}

func (o *modulo) String() string {
	return o.literal
}

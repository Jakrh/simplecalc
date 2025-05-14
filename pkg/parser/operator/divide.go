package operator

import "fmt"

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

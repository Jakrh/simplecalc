package operator

import "fmt"

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

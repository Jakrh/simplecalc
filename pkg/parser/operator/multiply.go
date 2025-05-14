package operator

import "fmt"

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

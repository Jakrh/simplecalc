package operator

import "fmt"
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

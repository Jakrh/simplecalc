package operator

import "fmt"

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

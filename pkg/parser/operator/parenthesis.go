package operator

import "fmt"

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

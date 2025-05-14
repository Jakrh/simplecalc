package parser_test

import (
	"errors"
	"math"
	"testing"

	"simplecalc/pkg/parser"
	"simplecalc/pkg/parser/operator"
)

func TestNewExpressionFromLexer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: nil,
		},
		{
			name:    "single number",
			input:   "42",
			want:    "42",
			wantErr: nil,
		},
		{
			name:    "simple addition",
			input:   "1 + 2",
			want:    "(+ 1 2)",
			wantErr: nil,
		},
		{
			name:    "simple subtraction",
			input:   "3 - 4",
			want:    "(- 3 4)",
			wantErr: nil,
		},
		{
			name:    "precedence",
			input:   "1 + 2 * 3",
			want:    "(+ 1 (* 2 3))",
			wantErr: nil,
		},
		{
			name:    "parentheses",
			input:   "(1 + 2) * 3",
			want:    "(* (+ 1 2) 3)",
			wantErr: nil,
		},
		{
			name:    "nested operations",
			input:   "4 * (5 + 6) - 7 / 8",
			want:    "(- (* 4 (+ 5 6)) (/ 7 8))",
			wantErr: nil,
		},
		{
			name:    "missing parentheses",
			input:   "(1 + 2",
			want:    "",
			wantErr: parser.ErrMissingRightParenthesis,
		},
		{
			name:    "negative number",
			input:   "-5",
			want:    "(- 0 5)",
			wantErr: nil,
		},
		{
			name:    "negative operation",
			input:   "-5 + 3",
			want:    "(+ (- 0 5) 3)",
			wantErr: nil,
		},
		{
			name:    "negative numbers with operators",
			input:   "-1 + -2 * -3",
			want:    "(+ (- 0 1) (* (- 0 2) (- 0 3)))",
			wantErr: nil,
		},
		{
			name:    "negative operation with parentheses",
			input:   "(-5 + 3) * 2",
			want:    "(* (+ (- 0 5) 3) 2)",
			wantErr: nil,
		},
		{
			name:    "negative leading decimal",
			input:   "-.5",
			want:    "(- 0 0.5)",
			wantErr: nil,
		},
		{
			name:    "negative leading decimal with operator",
			input:   "-.5 + 2",
			want:    "(+ (- 0 0.5) 2)",
			wantErr: nil,
		},
		{
			name:    "negative leading decimal with operator and parentheses",
			input:   "-5 * (2 + -.3)",
			want:    "(* (- 0 5) (+ 2 (- 0 0.3)))",
			wantErr: nil,
		},
		{
			name:    "positive leading decimal with operator and parentheses",
			input:   "-5 * (2 + +.3)",
			want:    "(* (- 0 5) (+ 2 (+ 0 0.3)))",
			wantErr: nil,
		},
		{
			name:    "negative leading parentheses",
			input:   "-(y + 2)",
			want:    "(- 0 (+ y 2))",
			wantErr: nil,
		},
		{
			name:    "leading right parentheses",
			input:   ")",
			want:    "",
			wantErr: parser.ErrMissingLeftParenthesis,
		},
		{
			name:    "missing left parentheses #1",
			input:   "1 + (2 * 3))",
			want:    "",
			wantErr: parser.ErrMissingLeftParenthesis,
		},
		{
			name:    "missing left parentheses #2",
			input:   "(1 * (2 + 3))) - 4",
			want:    "",
			wantErr: parser.ErrMissingLeftParenthesis,
		},
		{
			name:    "multiple letters variable assignment",
			input:   "var1 = 5",
			want:    "(= var1 5)",
			wantErr: nil,
		},
		{
			name:    "support negative variable",
			input:   "-y",
			want:    "(- 0 y)",
			wantErr: nil,
		},
		{
			name:    "support negative variable with parentheses",
			input:   "(-y)",
			want:    "(- 0 y)",
			wantErr: nil,
		},
		{
			name:    "support negative variable with parentheses and operator",
			input:   "(-y) + 2",
			want:    "(+ (- 0 y) 2)",
			wantErr: nil,
		},
		{
			name:    "support negative variable assignment",
			input:   "x = -y + 2",
			want:    "(= x (+ (- 0 y) 2))",
			wantErr: nil,
		},
		{
			name:    "simple power operation",
			input:   "2 ** 3",
			want:    "(** 2 3)",
			wantErr: nil,
		},
		{
			name:    "power operation with parentheses",
			input:   "(2 + 3) ** 2",
			want:    "(** (+ 2 3) 2)",
			wantErr: nil,
		},
		{
			name:    "power operation with negative number",
			input:   "-2 ** 3",
			want:    "(- 0 (** 2 3))",
			wantErr: nil,
		},
		{
			name:    "power operation with decimal number and parentheses",
			input:   "(2.5 + 3) ** 2",
			want:    "(** (+ 2.5 3) 2)",
			wantErr: nil,
		},
		{
			name:    "2 raised to the power of 0.5",
			input:   "2 ** 0.5",
			want:    "(** 2 0.5)",
			wantErr: nil,
		},
		{
			name:    "power operation with decimal number, negative sign and parentheses",
			input:   "-((2.5 * x) ** 6) ** y / .5 ** 3",
			want:    "(/ (- 0 (** (** (* 2.5 x) 6) y)) (** 0.5 3))",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer, err := parser.NewLexer(tt.input)
			if err != nil {
				t.Fatalf("NewLexer() error = %v", err)
				return
			}

			expr, err := parser.NewExpressionFromLexer(lexer)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewExpressionFromLexer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := expr.String()
			if got != tt.want {
				t.Errorf("Expression.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpressionEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      float64
		wantErr   error
		variables map[string]float64
	}{
		{
			name:    "empty input",
			input:   "",
			want:    0,
			wantErr: parser.ErrNilExpression,
		},
		{
			name:    "addition",
			input:   "1 + 2",
			want:    3,
			wantErr: nil,
		},
		{
			name:    "subtraction",
			input:   "5 - 3",
			want:    2,
			wantErr: nil,
		},
		{
			name:    "multiplication",
			input:   "4 * 2",
			want:    8,
			wantErr: nil,
		},
		{
			name:    "division",
			input:   "8 / 2",
			want:    4,
			wantErr: nil,
		},
		{
			name:    "division by zero",
			input:   "1 / 0",
			want:    0,
			wantErr: operator.ErrDivisionByZero,
		},
		{
			name:    "complex ops",
			input:   "2 + 3 * 4 - 6 / 2",
			want:    2 + 3*4 - 6.0/2.0,
			wantErr: nil,
		},
		{
			name:    "parentheses",
			input:   "(2 + 3) * (4 - 2)",
			want:    (2 + 3) * (4 - 2),
			wantErr: nil,
		},
		{
			name:    "nested",
			input:   "2 * (3 + (4 - 2))",
			want:    2 * (3 + (4 - 2)),
			wantErr: nil,
		},
		{
			name:    "negative number",
			input:   "-5 + 3",
			want:    -5 + 3,
			wantErr: nil,
		},
		{
			name:    "negative operation",
			input:   "-5 + 3 * 2",
			want:    -5 + 3*2,
			wantErr: nil,
		},
		{
			name:    "negative numbers with operators",
			input:   "-1 + -2 * -3",
			want:    -1 + -2*-3,
			wantErr: nil,
		},
		{
			name:    "negative operation with parentheses",
			input:   "(-5 + 3) * 2",
			want:    (-5 + 3) * 2,
			wantErr: nil,
		},
		{
			name:    "negative leading decimal",
			input:   "-.5",
			want:    -0.5,
			wantErr: nil,
		},
		{
			name:    "negative leading decimal with operator",
			input:   "-.5 + 2",
			want:    -0.5 + 2,
			wantErr: nil,
		},
		{
			name:    "negative leading decimal with operator and parentheses",
			input:   "-5 * (2 + -.3)",
			want:    -5 * (2 + -0.3),
			wantErr: nil,
		},
		{
			name:    "positive leading decimal with operator and parentheses",
			input:   "-5 * (2 + +.3)",
			want:    -5 * (2 + +0.3),
			wantErr: nil,
		},
		{
			name:      "negative leading parentheses",
			input:     "-(y + 2)",
			want:      -5,
			wantErr:   nil,
			variables: map[string]float64{"y": 3},
		},
		{
			name:    "parentheses only",
			input:   "()",
			want:    0,
			wantErr: parser.ErrNilExpression,
		},
		{
			name:      "simple variable operation",
			input:     "xx / yy",
			want:      5,
			wantErr:   nil,
			variables: map[string]float64{"xx": 10, "yy": 2},
		},
		{
			name:      "support negative variable",
			input:     "-y",
			want:      -12,
			wantErr:   nil,
			variables: map[string]float64{"y": 12},
		},
		{
			name:      "support negative variable with parentheses",
			input:     "(-y)",
			want:      -12,
			wantErr:   nil,
			variables: map[string]float64{"y": 12},
		},
		{
			name:      "support negative variable with parentheses and operator",
			input:     "(-y) + 2",
			want:      -10,
			wantErr:   nil,
			variables: map[string]float64{"y": 12},
		},
		{
			name:      "support negative variable assignment",
			input:     "x = -y + 2",
			want:      0,
			wantErr:   nil,
			variables: map[string]float64{"x": 7, "y": 12},
		},
		{
			// 1<<53 is 9007199254740992
			name:    "too large number",
			input:   "9007199254740992",
			want:    0,
			wantErr: parser.ErrNumOutOfRange,
		},
		{
			// -1<<53 is -9007199254740992
			name:    "too small number",
			input:   "-9007199254740992",
			want:    0,
			wantErr: parser.ErrNumOutOfRange,
		},
		{
			name:    "too large number from operation",
			input:   "9007199254740991 + 1",
			want:    0,
			wantErr: parser.ErrNumOutOfRange,
		},
		{
			name:    "too small number from operation",
			input:   "-9007199254740991 - 1",
			want:    0,
			wantErr: parser.ErrNumOutOfRange,
		},
		{
			name:      "too large number from variable",
			input:     "x",
			want:      0,
			wantErr:   parser.ErrNumOutOfRange,
			variables: map[string]float64{"x": 9007199254740992},
		},
		{
			name:      "too large number from variable with operation",
			input:     "x + 1",
			want:      0,
			wantErr:   parser.ErrNumOutOfRange,
			variables: map[string]float64{"x": 9007199254740991},
		},
		{
			name:    "simple power operation",
			input:   "2 ** 3",
			want:    8,
			wantErr: nil,
		},
		{
			name:    "power operation with parentheses",
			input:   "(2 + 3) ** 2",
			want:    25,
			wantErr: nil,
		},
		{
			name:    "power operation with negative number",
			input:   "-2 ** 3",
			want:    -8,
			wantErr: nil,
		},
		{
			name:    "power operation with decimal number and parentheses",
			input:   "(2.5 + 3) ** 2",
			want:    30.25,
			wantErr: nil,
		},
		{
			name:    "2 raised to the power of 0.5",
			input:   "2 ** 0.5",
			want:    1.4142135623730951,
			wantErr: nil,
		},
		{
			name:      "power operation with decimal number, negative sign and parentheses",
			input:     "-((2.5 * x) ** 6) ** y / .5 ** 3",
			want:      -64.0,
			wantErr:   nil,
			variables: map[string]float64{"x": 1.6, "y": 0.25},
		},
		{
			name:    "trailing operator",
			input:   "2 + ",
			want:    0,
			wantErr: operator.ErrInvalidOperandCount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer, err := parser.NewLexer(tt.input)
			if err != nil {
				t.Fatalf("NewLexer() error = %v", err)
				return
			}

			expr, err := parser.NewExpressionFromLexer(lexer)
			if err != nil {
				t.Fatalf("NewExpressionFromLexer() error = %v", err)
				return
			}

			if tt.variables == nil {
				tt.variables = make(map[string]float64)
			}
			got, err := expr.Evaluate(tt.variables)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if math.Abs(got-tt.want) > parser.IntApproxTolerance {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

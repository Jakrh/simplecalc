package parser_test

import (
	"errors"
	"slices"
	"testing"

	"simplecalc/pkg/parser"
	"simplecalc/pkg/parser/operator"
)

func TestParser_Parse_MultiLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []float64
		wantErr error
	}{
		{
			name:  "simple multi-line",
			input: "1 + 1; 2 * 3",
			want:  []float64{2, 6},
		},
		{
			name:  "multi-line with assignment",
			input: "a = 5; a * 2",
			want:  []float64{10},
		},
		{
			name:  "multi-line with multiple assignments",
			input: "a = 2; b = 3; a + b",
			want:  []float64{5},
		},
		{
			name:  "multi-line assignment only, last one is assignment",
			input: "a = 5; b = a * 2",
			want:  []float64{},
		},
		{
			name:  "multi-line assignment only, single assignment",
			input: "a = 5",
			want:  []float64{},
		},
		{
			name:  "empty input",
			input: "",
			want:  []float64{},
		},
		{
			name:  "only semicolons",
			input: ";;;",
			want:  []float64{},
		},
		{
			name:  "semicolon at the end",
			input: "1+1;",
			want:  []float64{2},
		},
		{
			name:  "semicolon at the beginning",
			input: ";1+1",
			want:  []float64{2},
		},
		{
			name:  "multiple expressions with spaces and semicolons",
			input: "  x = 10;  y = x / 2 ;  y + 3 ",
			want:  []float64{8},
		},
		{
			name:  "assignment followed by expression using the variable",
			input: "var = 15; var - 5",
			want:  []float64{10},
		},
		{
			name:    "error in later statement",
			input:   "a = 1; 1 / 0",
			wantErr: operator.ErrDivisionByZero,
		},
		{
			name:    "error in first statement",
			input:   "1 / 0; a = 1",
			wantErr: operator.ErrDivisionByZero,
		},
		{
			name:  "complex multi-line with assignments and evaluation",
			input: "a=1;b=2;c=a+b;c*2",
			want:  []float64{6},
		},
		{
			name:  "complex multi-line with assignments only",
			input: "a=1;b=2;c=a+b",
			want:  []float64{},
		},
		{
			name:  "multi-line with final statement being an expression",
			input: "x=100; y=200; x+y",
			want:  []float64{300},
		},
		{
			name:  "multi-line with final statement being an assignment",
			input: "x=100; y=200; z=x+y",
			want:  []float64{},
		},
		{
			name:  "multi-line with empty statements",
			input: "a=1;;b=a+1;;b*2",
			want:  []float64{4},
		},
		{
			name:  "multi-line with leading/trailing spaces and empty statements",
			input: "  a = 5 ; ; b = a + 5 ;  b * 2  ",
			want:  []float64{20},
		},
		{
			name:  "round if approximate to integer",
			input: "1.99999999999",
			want:  []float64{2},
		},
		{
			name:  "round negative number if approximate to integer",
			input: "-1.99999999999",
			want:  []float64{-2},
		},
		{
			name:  "round negative number with a small decimal if approximate to integer",
			input: "-2.00000000001",
			want:  []float64{-2},
		},
		{
			name:  "not round if not approximate to integer",
			input: "1.9999999999",
			want:  []float64{1.9999999999},
		},
		{
			name:  "not round negative number if not approximate to integer",
			input: "-1.9999999999",
			want:  []float64{-1.9999999999},
		},
		{
			name:  "2 raised to the power of 0.5 with assignment",
			input: "x = 2; y = 0.5; z = x ** y; z",
			want:  []float64{1.4142135623730951},
		},
		{
			name:    "power operation with decimal number, negative sign and parentheses",
			input:   "x = 1.6; y = .25; -((2.5 * x) ** 6) ** y / .5 ** 3",
			want:    []float64{-64.0},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser()
			got, err := p.Parse(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Parse() error = %q, wantErr %q", err, tt.wantErr)
				return
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

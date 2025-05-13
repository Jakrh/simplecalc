package parser_test

import (
	"reflect"
	"testing"

	"simplecalc/pkg/parser"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []parser.Token
		wantErr bool
	}{
		{
			name:  "empty input",
			input: "",
			want: []parser.Token{
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "simple addition",
			input: "1 + 2",
			want: []parser.Token{
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "whitespace and multicharacters",
			input: " 12 - 3.4 ",
			want: []parser.Token{
				parser.NewAtomNumToken("12"),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("3.4"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:    "illegal char",
			input:   "@",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "multiple operators",
			input: "1 + 2 * 3 - 4 / 5",
			want: []parser.Token{
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewAtomNumToken("3"),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("4"),
				parser.NewOPToken(parser.TokenDivide),
				parser.NewAtomNumToken("5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "parentheses",
			input: "(1 + 2) * 3",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewAtomNumToken("3"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "integer",
			input: "123",
			want: []parser.Token{
				parser.NewAtomNumToken("123"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "decimal",
			input: "3.14",
			want: []parser.Token{
				parser.NewAtomNumToken("3.14"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "leading decimal",
			input: ".5",
			want: []parser.Token{
				parser.NewAtomNumToken(".5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "trailing decimal",
			input: "10.",
			want: []parser.Token{
				parser.NewAtomNumToken("10."),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:    "double decimal",
			input:   "1.2.3",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "extra parentheses",
			input: "1 + 2)",
			want: []parser.Token{
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative number",
			input: "-5",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative decimal",
			input: "-3.14",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("3.14"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative trailing decimal",
			input: "-10.",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("10."),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative numbers with operators",
			input: "-1 + -2 * -3",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("3"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative number with parentheses",
			input: "(-5 + 3) * 2",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("5"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("3"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewAtomNumToken("2"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative number minus operator",
			input: "-2 - 3",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("3"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal",
			input: "-.5",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken(".5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal with operator",
			input: "-.5 + 2",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken(".5"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal with operator and parentheses",
			input: "-5 * (2 + -.3)",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("5"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken(".3"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "positive leading decimal with operator and parentheses",
			input: "-5 * (2 + +.3)",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("5"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken(".3"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "leading right parentheses",
			input: ")",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "missing left parentheses",
			input: "1 + (2 * 3))",
			want: []parser.Token{
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewAtomNumToken("3"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "missing left parentheses with tailing operation",
			input: "(1 * (2 + 3))) - 4",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("1"),
				parser.NewOPToken(parser.TokenMultiply),
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewAtomNumToken("2"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("3"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("4"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "parentheses only",
			input: "()",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "multiple letters variable assignment",
			input: "var1 = 5",
			want: []parser.Token{
				parser.NewAtomVarToken("var1"),
				parser.NewOPToken(parser.TokenAssign),
				parser.NewAtomNumToken("5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "variable assignment with a negative number",
			input: "x = -5",
			want: []parser.Token{
				parser.NewAtomVarToken("x"),
				parser.NewOPToken(parser.TokenAssign),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomNumToken("5"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "support negative variable",
			input: "-y",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomVarToken("y"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "support negative variable with parentheses",
			input: "(-y)",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomVarToken("y"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "support negative variable with parentheses and operator",
			input: "(-y) + 2",
			want: []parser.Token{
				parser.NewOPToken(parser.TokenLeftParen),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomVarToken("y"),
				parser.NewOPToken(parser.TokenRightParen),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
		{
			name:  "support negative variable assignment",
			input: "x = -y + 2",
			want: []parser.Token{
				parser.NewAtomVarToken("x"),
				parser.NewOPToken(parser.TokenAssign),
				parser.NewOPToken(parser.TokenMinus),
				parser.NewAtomVarToken("y"),
				parser.NewOPToken(parser.TokenPlus),
				parser.NewAtomNumToken("2"),
				parser.NewEOFToken(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := parser.NewLexer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLexer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			var got []parser.Token
			for l.HasNext() {
				got = append(got, l.Next())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens = %v, want %v", got, tt.want)
			}
		})
	}
}

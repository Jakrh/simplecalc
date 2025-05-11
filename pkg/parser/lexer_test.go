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
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "simple addition",
			input: "1 + 2",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "whitespace and multicharacters",
			input: " 12 - 3.4 ",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "12"},
				{Type: parser.TokenMinus, Literal: "-"},
				{Type: parser.TokenAtom, Literal: "3.4"},
				{Type: parser.TokenEOF, Literal: ""},
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
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenMinus, Literal: "-"},
				{Type: parser.TokenAtom, Literal: "4"},
				{Type: parser.TokenDivide, Literal: "/"},
				{Type: parser.TokenAtom, Literal: "5"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "parentheses",
			input: "(1 + 2) * 3",
			want: []parser.Token{
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "integer",
			input: "123",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "123"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "decimal",
			input: "3.14",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "3.14"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "leading decimal",
			input: ".5",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: ".5"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "trailing decimal",
			input: "10.",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "10."},
				{Type: parser.TokenEOF, Literal: ""},
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
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative number",
			input: "-5",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-5"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative decimal",
			input: "-3.14",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-3.14"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative trailing decimal",
			input: "-10.",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-10."},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative numbers with operators",
			input: "-1 + -2 * -3",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "-2"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenAtom, Literal: "-3"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative number with parentheses",
			input: "(-5 + 3) * 2",
			want: []parser.Token{
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "-5"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative number minus operator",
			input: "-2 - 3",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-2"},
				{Type: parser.TokenMinus, Literal: "-"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal",
			input: "-.5",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-.5"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal with operator",
			input: "-.5 + 2",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-.5"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "negative leading decimal with operator and parentheses",
			input: "-5 * (2 + -.3)",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "-5"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "-.3"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "leading right parentheses",
			input: ")",
			want: []parser.Token{
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "missing left parentheses",
			input: "1 + (2 * 3))",
			want: []parser.Token{
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "missing left parentheses with tailing operation",
			input: "(1 * (2 + 3))) - 4",
			want: []parser.Token{
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "1"},
				{Type: parser.TokenMultiply, Literal: "*"},
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenAtom, Literal: "2"},
				{Type: parser.TokenPlus, Literal: "+"},
				{Type: parser.TokenAtom, Literal: "3"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenMinus, Literal: "-"},
				{Type: parser.TokenAtom, Literal: "4"},
				{Type: parser.TokenEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "parentheses only",
			input: "()",
			want: []parser.Token{
				{Type: parser.TokenLeftParen, Literal: "("},
				{Type: parser.TokenRightParen, Literal: ")"},
				{Type: parser.TokenEOF, Literal: ""},
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

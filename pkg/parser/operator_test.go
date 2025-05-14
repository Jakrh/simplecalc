package parser_test

import (
	"testing"

	"simplecalc/pkg/parser"
)

func TestLexWithOperator(t *testing.T) {
	type input struct {
		input  string
		cursor int
	}
	type output struct {
		parser.Operator
		newCursor int
	}
	tests := []struct {
		name  string
		input input
		want  output
	}{
		{
			name: "handle plus operator from simple expression",
			input: input{
				input:  "1+2",
				cursor: 1,
			},
			want: output{
				Operator:  parser.GetOperator("+"),
				newCursor: 2,
			},
		},
		{
			name: "handle minus operator from simple expression",
			input: input{
				input:  "1-2",
				cursor: 1,
			},
			want: output{
				Operator:  parser.GetOperator("-"),
				newCursor: 2,
			},
		},
		{
			name: "handle multiply operator from simple expression",
			input: input{
				input:  "1*2",
				cursor: 1,
			},
			want: output{
				Operator:  parser.GetOperator("*"),
				newCursor: 2,
			},
		},
		{
			name: "handle divide operator from simple expression",
			input: input{
				input:  "1/2",
				cursor: 1,
			},
			want: output{
				Operator:  parser.GetOperator("/"),
				newCursor: 2,
			},
		},
		{
			name: "handle power operator from simple expression",
			input: input{
				input:  "1**2",
				cursor: 1,
			},
			want: output{
				Operator:  parser.GetOperator("**"),
				newCursor: 3,
			},
		},
		{
			name: "handle number from single digit",
			input: input{
				input:  "1",
				cursor: 0,
			},
			want: output{
				Operator:  nil,
				newCursor: 1,
			},
		},
		{
			name: "handle variable from single letter",
			input: input{
				input:  "a",
				cursor: 0,
			},
			want: output{
				Operator:  nil,
				newCursor: 1,
			},
		},
		{
			name: "handle number from a expression",
			input: input{
				input:  "1+2",
				cursor: 2,
			},
			want: output{
				Operator:  nil,
				newCursor: 3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op, cursor := parser.LexWithOperator(&test.input.input, test.input.cursor)
			if op != test.want.Operator {
				t.Errorf("got operator '%s', want '%s'", op, test.want.Operator)
			}
			if cursor != test.want.newCursor {
				t.Errorf("got cursor %d, want %d", cursor, test.want.newCursor)
			}
		})
	}
}

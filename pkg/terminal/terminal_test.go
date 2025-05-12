package terminal

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"golang.org/x/term"
)

func TestTerminal_ReadLine(t *testing.T) {
	// inputs must be end with CR or CRLF
	// because ReadLine() treats CR
	// as line ending, not LF.

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "normal input",
			input:   "hello world\r",
			want:    "hello world",
			wantErr: nil,
		},
		{
			name:    "empty input",
			input:   "\r",
			want:    "",
			wantErr: nil,
		},
		{
			name:    "EOF error",
			input:   "",
			want:    "",
			wantErr: io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare in/out buffers
			var buf bytes.Buffer
			buf.WriteString(tt.input)
			termIn := &buf
			termOut := &bytes.Buffer{}

			// create a terminal.Terminal to use underlying ReadLine
			tTerm := term.NewTerminal(struct {
				io.Reader
				io.Writer
			}{termIn, termOut}, "> ")
			trm := &Terminal{terminal: tTerm}

			// read line from terminal and check result
			line, err := trm.ReadLine()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ReadLine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if line != tt.want {
				t.Errorf("ReadLine() = %q, want %q", line, tt.want)
			}
		})
	}
}

func TestTerminal_GetAndClearHistory(t *testing.T) {
	trm := &Terminal{history: &History{}}

	// history is empty initially
	if got := trm.GetHistory(); got != "" {
		t.Errorf("GetHistory() = %q, want empty string", got)
	}

	// add entries via history directly
	trm.history.Add("one")
	trm.history.Add("two")
	want := "one\r\ntwo"
	if got := trm.GetHistory(); got != want {
		t.Errorf("GetHistory() = %q, want %q", got, want)
	}

	// clear and verify
	trm.ClearHistory()
	if got := trm.GetHistory(); got != "" {
		t.Errorf("GetHistory() after ClearHistory = %q, want empty string", got)
	}
}

func TestTerminal_Restore(t *testing.T) {
	// recover from panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Restore() panic = %v", r)
		}
	}()

	trm := &Terminal{}

	// test restore when oldState is nil
	trm.oldState = nil
	if err := trm.Restore(); err == nil {
		t.Errorf("Restore() error = %v, want error", err)
	}
}

func TestNewTerminal_Error(t *testing.T) {
	// use a non-terminal file (e.g., /dev/null) so MakeRaw fails
	f, err := os.OpenFile(os.DevNull, os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("os.OpenFile failed: %v", err)
	}
	defer f.Close()

	// try to create a new terminal with a non-terminal file
	trm, err := NewTerminal(f, "prompt")
	if err == nil {
		t.Fatalf("NewTerminal() expected error for non-terminal file, got nil, trm=%v", trm)
	}
}

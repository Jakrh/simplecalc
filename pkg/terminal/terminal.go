package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var (
	ErrCreateTerminal = fmt.Errorf("error creating terminal")
	ErrNoOldState     = fmt.Errorf("no old state to restore")
	ErrNoInputFile    = fmt.Errorf("no input file to restore")
	ErrSetRawMode     = fmt.Errorf("error setting terminal to raw mode")
)

type Terminal struct {
	prompt    string
	oldState  *term.State
	terminal  *term.Terminal
	history   *History
	inputFile *os.File
}

func NewTerminal(f *os.File, prompt string) (*Terminal, error) {
	// Set the terminal to raw mode
	fd := int(f.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// May return error as syscall.Errno
		return nil, fmt.Errorf("error setting terminal to raw mode: %w", err)
	}

	// Create a new terminal instance
	t := term.NewTerminal(f, prompt)
	if t == nil {
		return nil, ErrCreateTerminal
	}

	// Set the history instance
	h := &History{
		history: make([]string, 0),
	}
	t.History = h

	return &Terminal{
		prompt:    prompt,
		oldState:  oldState,
		terminal:  t,
		history:   h,
		inputFile: f,
	}, nil
}

// ReadLine wraps the ReadLine method of the term
func (t *Terminal) ReadLine() (string, error) {
	line, err := t.terminal.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

// Restore wraps the Restore method of the term
func (t *Terminal) Restore() error {
	if t.oldState == nil {
		return ErrNoOldState
	}
	if t.inputFile == nil {
		return ErrNoInputFile
	}

	if err := term.Restore(int(t.inputFile.Fd()), t.oldState); err != nil {
		return err
	}
	return nil
}

func (t *Terminal) GetHistory() string {
	if t.history == nil {
		return ""
	}

	return t.history.String()
}

func (t *Terminal) ClearHistory() {
	if t.history == nil {
		return
	}

	t.history.Clear()
}

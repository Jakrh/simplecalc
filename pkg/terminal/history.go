package terminal

import (
	"fmt"
	"strings"
)

type History struct {
	history []string
}

func (h *History) Add(entry string) {
	// Not to add empty entries to history
	if entry == "" {
		return
	}

	h.history = append(h.history, entry)
}

func (h *History) Len() int {
	return len(h.history)
}

func (h *History) At(idx int) string {
	if idx < 0 || idx >= len(h.history) {
		panic(fmt.Sprintf("index out of range: %d", idx))
	}

	return h.history[len(h.history)-idx-1]
}

func (h *History) Clear() {
	h.history = []string{}
}

func (h *History) String() string {
	if len(h.history) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, line := range h.history {
		sb.WriteString(line)
		if i != len(h.history)-1 {
			sb.WriteString("\r\n")
		}
	}

	return sb.String()
}

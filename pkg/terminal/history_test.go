package terminal

import (
	"testing"
)

func TestHistory_AddAndLen(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
		wantLen int
	}{
		{
			name:    "empty",
			entries: []string{},
			wantLen: 0,
		},
		{
			name:    "single",
			entries: []string{"a"},
			wantLen: 1,
		},
		{
			name:    "multiple",
			entries: []string{"a", "b", "c"},
			wantLen: 3,
		},
		{
			name:    "multiple with empty",
			entries: []string{"a", "", "b"}, // the empty one should be ignored
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h History
			for _, e := range tt.entries {
				h.Add(e)
			}
			if got := h.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
		})
	}
}

func TestHistory_At(t *testing.T) {
	entries := []string{"one", "two", "three"}
	var h History
	for _, e := range entries {
		h.Add(e)
	}

	tests := []struct {
		name      string
		idx       int
		want      string
		wantPanic bool
	}{
		{
			name:      "first",
			idx:       0,
			want:      "three",
			wantPanic: false,
		},
		{
			name:      "second",
			idx:       1,
			want:      "two",
			wantPanic: false,
		},
		{
			name:      "third",
			idx:       2,
			want:      "one",
			wantPanic: false,
		},
		{
			name:      "neg",
			idx:       -1,
			want:      "",
			wantPanic: true,
		},
		{
			name:      "out of range",
			idx:       3,
			want:      "",
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("At(%d) panic = %v, want panic %v", tt.idx, r != nil, tt.wantPanic)
				}
			}()
			got := h.At(tt.idx)
			if !tt.wantPanic && got != tt.want {
				t.Errorf("At(%d) = %q, want %q", tt.idx, got, tt.want)
			}
		})
	}
}

func TestHistory_Clear(t *testing.T) {
	var h History
	h.Add("test")
	if h.Len() != 1 {
		t.Fatalf("Len before clear = %d, want 1", h.Len())
	}
	h.Clear()
	if h.Len() != 0 {
		t.Errorf("Len after clear = %d, want 0", h.Len())
	}
}

func TestHistory_String(t *testing.T) {
	var h History
	if got := h.String(); got != "" {
		t.Errorf("String() empty history = %q, want empty", got)
	}
	entries := []string{"first", "second", "third"}
	for _, e := range entries {
		h.Add(e)
	}
	wantStr := "first\r\nsecond\r\nthird"
	if got := h.String(); got != wantStr {
		t.Errorf("String() = %q, want %q", got, wantStr)
	}
}

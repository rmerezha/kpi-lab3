package lang

import (
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		line        string
		expectError bool
	}{
		{"white", false},
		{"green", false},
		{"update", false},
		{"reset", false},
		{"bgrect 0.1 0.2 0.3 0.4", false},
		{"figure 0.5 0.6", false},
		{"move 0.1 0.2", false},
		{"unknown", true},
		{"white extra", true},
		{"bgrect 1.0", true},
	}

	for _, tt := range tests {
		p := &Parser{}
		err := p.parseCommand(tt.line)
		if tt.expectError && err == nil {
			t.Errorf("expected error for %q, got nil", tt.line)
		}
		if !tt.expectError && err != nil {
			t.Errorf("unexpected error for %q: %v", tt.line, err)
		}
	}
}

func TestParse(t *testing.T) {
	input := `
white
bgrect 0.1 0.1 0.2 0.2
figure 0.5 0.5
move 0.1 0.1
update
`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(ops) == 0 {
		t.Errorf("expected operations, got 0")
	}
}

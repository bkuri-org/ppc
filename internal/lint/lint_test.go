package lint

import (
	"testing"

	"github.com/bkuri/ppc/internal/model"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello world", 2},
		{"one two three four", 4},
		{"", 0},
		{"   ", 0},
		{"word", 1},
		{"multiple   spaces   between", 3},
	}

	for _, tc := range tests {
		got := countWords(tc.input)
		if got != tc.expected {
			t.Errorf("countWords(%q) = %d, want %d", tc.input, got, tc.expected)
		}
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"single line", 1},
		{"two\nlines", 2},
		{"three\nlines\nhere", 3},
		{"trailing\n", 2},
	}

	for _, tc := range tests {
		got := countLines(tc.input)
		if got != tc.expected {
			t.Errorf("countLines(%q) = %d, want %d", tc.input, got, tc.expected)
		}
	}
}

func TestPercentOver(t *testing.T) {
	tests := []struct {
		actual    int
		threshold int
		expected  int
	}{
		{150, 100, 50},
		{200, 100, 100},
		{105, 100, 5},
		{100, 100, 0},
		{50, 100, -50},
		{100, 0, 0},
	}

	for _, tc := range tests {
		got := percentOver(tc.actual, tc.threshold)
		if got != tc.expected {
			t.Errorf("percentOver(%d, %d) = %d, want %d", tc.actual, tc.threshold, got, tc.expected)
		}
	}
}

func TestTagPatternMatches(t *testing.T) {
	tags := []string{"risk:low", "domain:api", "status:active"}

	tests := []struct {
		pattern  string
		expected bool
	}{
		{"risk:low", true},
		{"risk:high", false},
		{"risk:*", true},
		{"domain:*", true},
		{"tier:*", false},
		{"status:active", true},
	}

	for _, tc := range tests {
		got := tagPatternMatches(tc.pattern, tags)
		if got != tc.expected {
			t.Errorf("tagPatternMatches(%q, %v) = %v, want %v", tc.pattern, tags, got, tc.expected)
		}
	}
}

func TestHasField(t *testing.T) {
	tests := []struct {
		field    string
		expected bool
	}{
		{"id", true},
		{"desc", true},
		{"priority", false},
		{"tags", false},
		{"requires", false},
		{"unknown", false},
	}

	fm := model.Frontmatter{
		ID:   "test",
		Desc: "Test module",
	}

	for _, tc := range tests {
		got := hasField(fm, tc.field)
		if got != tc.expected {
			t.Errorf("hasField(%+v, %q) = %v, want %v", fm, tc.field, got, tc.expected)
		}
	}
}

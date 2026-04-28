package model

import "testing"

func TestLayerName(t *testing.T) {
	tests := []struct {
		idx      int
		expected string
	}{
		{0, "base"},
		{1, "modes"},
		{2, "traits"},
		{3, "policies"},
		{4, "contracts"},
		{5, "guardrails"},
		{-1, "unknown"},
		{6, "unknown"},
		{99, "unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			got := LayerName(tc.idx)
			if got != tc.expected {
				t.Errorf("LayerName(%d) = %q, want %q", tc.idx, got, tc.expected)
			}
		})
	}
}

func TestLayerIndexFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"base.md", 0},
		{"base/content.md", 0},
		{"modes/explore.md", 1},
		{"traits/terse.md", 2},
		{"policies/rule.md", 3},
		{"contracts/api.md", 4},
		{"guardrails/safety.md", 5},
		{"unknown/foo.md", 0},   // fallback to 0
		{"foo/modes/bar.md", 1}, // modes found in path
		{"", 0},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := LayerIndexFromPath(tc.path)
			if got != tc.expected {
				t.Errorf("LayerIndexFromPath(%q) = %d, want %d", tc.path, got, tc.expected)
			}
		})
	}
}

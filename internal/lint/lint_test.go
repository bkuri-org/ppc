package lint

import (
	"testing"

	"github.com/bkuri/ppc/internal/loader"
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

func TestRunBasicStats(t *testing.T) {
	result, err := Run("testdata", Config{})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Stats["module_count"] != 6 {
		t.Errorf("module_count = %d, want 6", result.Stats["module_count"])
	}

	if result.Stats["word_count"] == 0 {
		t.Error("word_count should be > 0")
	}

	if result.Stats["line_count"] == 0 {
		t.Error("line_count should be > 0")
	}
}

func TestRunMaxWords(t *testing.T) {
	result, err := Run("testdata", Config{MaxWords: 10})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "max_words" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected max_words violation")
	}
}

func TestRunMaxLines(t *testing.T) {
	result, err := Run("testdata", Config{MaxLines: 2})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "max_lines" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected max_lines violation")
	}
}

func TestRunMaxModules(t *testing.T) {
	result, err := Run("testdata", Config{MaxModules: 3})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "max_modules" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected max_modules violation")
	}
}

func TestRunForbidEmptyBody(t *testing.T) {
	result, err := Run("testdata", Config{ForbidEmptyBody: true})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "forbid_empty_body" && v.Module == "traits/empty" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected forbid_empty_body violation for traits/empty")
	}
}

func TestRunForbidTags(t *testing.T) {
	result, err := Run("testdata", Config{ForbidTags: []string{"risk:low"}})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "forbid_tags" && v.Module == "modes/test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected forbid_tags violation for modes/test")
	}
}

func TestRunRequireTags(t *testing.T) {
	result, err := Run("testdata", Config{RequireTags: []string{"nonexistent:*"}})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "require_tags" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected require_tags violation")
	}
}

func TestRunRequireTagsMet(t *testing.T) {
	result, err := Run("testdata", Config{RequireTags: []string{"risk:*"}})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, v := range result.Violations {
		if v.Rule == "require_tags" {
			t.Error("unexpected require_tags violation")
		}
	}
}

func TestRunMaxModuleWords(t *testing.T) {
	result, err := Run("testdata", Config{MaxModuleWords: 3})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "max_module_words" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected max_module_words violation")
	}
}

func TestCalculateModuleDepth(t *testing.T) {
	tests := []struct {
		id          string
		minDepth    int
		chainPrefix string
	}{
		{"base", 1, "base"},
		{"traits/deep3", 2, "traits/deep3"},
		{"traits/deep2", 3, "traits/deep2"},
		{"traits/deep1", 4, "traits/deep1"},
	}

	modByID, _ := loader.LoadModules("testdata")

	for _, tc := range tests {
		depth, chain := calculateModuleDepth(modByID, tc.id)
		if depth < tc.minDepth {
			t.Errorf("depth(%s) = %d, want >= %d", tc.id, depth, tc.minDepth)
		}
		if len(chain) == 0 || chain[0] != tc.chainPrefix {
			t.Errorf("chain(%s) = %v, want starting with %s", tc.id, chain, tc.chainPrefix)
		}
	}
}

func TestRunMaxDepth(t *testing.T) {
	result, err := Run("testdata", Config{MaxDepth: 2})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, v := range result.Violations {
		if v.Rule == "max_depth" {
			found = true
			if v.Module == "traits/deep1" {
				if !contains(v.Message, "chain =") {
					t.Error("max_depth violation should include chain")
				}
			}
		}
	}
	if !found {
		t.Error("expected max_depth violation")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

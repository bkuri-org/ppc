package compile

import (
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	t.Run("valid compilation", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "explore",
			Contract:   "simple",
			PromptsDir: "testdata",
		}

		out, meta, err := Compile(opts)
		if err != nil {
			t.Fatalf("Compile failed: %s", err)
		}

		if out == "" {
			t.Fatal("expected non-empty output")
		}

		// Output should contain module bodies
		if !strings.Contains(out, "helpful assistant") {
			t.Error("output should contain base module content")
		}

		// Meta should have selected IDs
		if len(meta.SelectedIDs) == 0 {
			t.Error("meta.SelectedIDs should not be empty")
		}

		// Hash should be a non-empty SHA256 hex string
		if len(meta.Hash) != 64 {
			t.Errorf("meta.Hash length = %d, want 64", len(meta.Hash))
		}

		// Order should match closure
		if len(meta.Order) != len(meta.ClosureIDs) {
			t.Errorf("Order length %d != ClosureIDs length %d",
				len(meta.Order), len(meta.ClosureIDs))
		}
	})

	t.Run("nonexistent prompts directory", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "explore",
			Contract:   "simple",
			PromptsDir: "testdata/nonexistent",
		}

		_, _, err := Compile(opts)
		if err == nil {
			t.Fatal("expected error for nonexistent prompts directory")
		}
	})

	t.Run("missing mode module", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "nonexistent_mode",
			Contract:   "simple",
			PromptsDir: "testdata",
		}

		_, _, err := Compile(opts)
		if err == nil {
			t.Fatal("expected error for missing mode module")
		}
	})

	t.Run("missing contract module", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "explore",
			Contract:   "nonexistent_contract",
			PromptsDir: "testdata",
		}

		_, _, err := Compile(opts)
		if err == nil {
			t.Fatal("expected error for missing contract module")
		}
	})

	t.Run("deterministic output", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "explore",
			Contract:   "simple",
			PromptsDir: "testdata",
		}

		out1, _, err1 := Compile(opts)
		if err1 != nil {
			t.Fatalf("first compile failed: %s", err1)
		}

		out2, _, err2 := Compile(opts)
		if err2 != nil {
			t.Fatalf("second compile failed: %s", err2)
		}

		if out1 != out2 {
			t.Error("compilation should be deterministic")
		}
	})

	t.Run("compilation with vars", func(t *testing.T) {
		opts := CompileOptions{
			Mode:       "explore",
			Contract:   "simple",
			PromptsDir: "testdata",
			Vars:       map[string]any{"name": "test-value"},
		}

		out, _, err := Compile(opts)
		if err != nil {
			t.Fatalf("Compile with vars failed: %s", err)
		}
		if out == "" {
			t.Fatal("expected non-empty output with vars")
		}
	})

	t.Run("buildSelectedIDs", func(t *testing.T) {
		opts := CompileOptions{
			Mode:     "explore",
			Contract: "simple",
			Traits:   []string{"traits/terse"},
		}

		ids := buildSelectedIDs(opts)
		expected := []string{"base", "modes/explore", "contracts/simple", "traits/terse"}
		if len(ids) != len(expected) {
			t.Fatalf("buildSelectedIDs length = %d, want %d", len(ids), len(expected))
		}
		for i, id := range ids {
			if id != expected[i] {
				t.Errorf("ids[%d] = %q, want %q", i, id, expected[i])
			}
		}
	})
}

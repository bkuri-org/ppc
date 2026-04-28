package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bkuri/ppc/internal/compile"
)

func TestGoldenSnapshots(t *testing.T) {
	tests := []struct {
		name    string
		opts    compile.CompileOptions
		fixture string
	}{
		{
			name: "explore_conservative_revisions1",
			opts: compile.CompileOptions{
				Mode:       "explore",
				Contract:   "markdown",
				Traits:     []string{"traits/conservative"},
				PromptsDir: filepath.Join("..", "prompts"),
				Vars:       map[string]any{"mode": "explore", "revisions": "1"},
			},
			fixture: "testdata/explore_conservative_revisions1.md",
		},
		{
			name: "build_creative_code",
			opts: compile.CompileOptions{
				Mode:       "build",
				Contract:   "code",
				Traits:     []string{"traits/creative"},
				PromptsDir: filepath.Join("..", "prompts"),
				Vars:       map[string]any{"mode": "build"},
			},
			fixture: "testdata/build_creative_code.md",
		},
		{
			name: "ship_conservative_terse",
			opts: compile.CompileOptions{
				Mode:       "ship",
				Contract:   "markdown",
				Traits:     []string{"traits/conservative", "traits/terse"},
				PromptsDir: filepath.Join("..", "prompts"),
				Vars:       map[string]any{"mode": "ship"},
			},
			fixture: "testdata/ship_conservative_terse.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, _, err := compile.Compile(tt.opts)
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			expected, err := os.ReadFile(tt.fixture)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			// Byte-for-byte comparison
			if string(expected) != out {
				t.Errorf("output differs from fixture")
				t.Logf("expected (%d bytes):\n%s", len(expected), string(expected[:min(len(expected), 200)]))
				t.Logf("got (%d bytes):\n%s", len(out), out[:min(len(out), 200)])
			}
		})
	}
}

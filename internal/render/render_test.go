package render

import (
	"strings"
	"testing"

	"github.com/bkuri/ppc/internal/model"
	"github.com/bkuri/ppc/internal/substitute"
)

func TestRender(t *testing.T) {
	t.Run("empty module list", func(t *testing.T) {
		out, unresolved := Render(nil, nil)
		// Render always appends a trailing newline
		if out != "\n" {
			t.Errorf("Render(nil, nil) = %q, want %q", out, "\n")
		}
		if len(unresolved) != 0 {
			t.Errorf("expected no unresolved vars, got %v", unresolved)
		}
	})

	t.Run("single module", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "Hello world."},
		}
		out, _ := Render(mods, nil)
		if !strings.Contains(out, "Hello world.") {
			t.Errorf("output missing body content: %q", out)
		}
		if !strings.HasSuffix(out, "\n") {
			t.Error("output should end with newline")
		}
	})

	t.Run("multiple modules separated by blank lines", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "Base content."},
			{Front: model.Frontmatter{ID: "modes/explore"}, Body: "Mode content."},
			{Front: model.Frontmatter{ID: "traits/terse"}, Body: "Trait content."},
		}
		out, _ := Render(mods, nil)

		if !strings.Contains(out, "Base content.") {
			t.Error("output missing base content")
		}
		if !strings.Contains(out, "Mode content.") {
			t.Error("output missing mode content")
		}
		if !strings.Contains(out, "Trait content.") {
			t.Error("output missing trait content")
		}

		// Check modules are separated by double newline
		if !strings.Contains(out, "Base content.\n\nMode content.") {
			t.Error("modules should be separated by blank line")
		}
	})

	t.Run("variable substitution", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "Hello {{name}}."},
		}
		vars := substitute.Vars{"name": "world"}
		out, unresolved := Render(mods, vars)

		if !strings.Contains(out, "Hello world.") {
			t.Errorf("expected variable substitution, got: %q", out)
		}
		if len(unresolved) != 0 {
			t.Errorf("expected no unresolved vars, got %v", unresolved)
		}
	})

	t.Run("unresolved variables tracked", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "Hello {{name}} and {{missing}}."},
		}
		out, unresolved := Render(mods, substitute.Vars{"name": "world"})

		if !strings.Contains(out, "Hello world and {{missing}}.") {
			t.Errorf("resolved var should be substituted, got: %q", out)
		}
		if len(unresolved) != 1 || unresolved[0] != "missing" {
			t.Errorf("expected [missing] unresolved, got %v", unresolved)
		}
	})

	t.Run("trailing newlines stripped from bodies", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "Content.\n\n"},
			{Front: model.Frontmatter{ID: "traits/a"}, Body: "More.\n"},
		}
		out, _ := Render(mods, nil)

		// Output should have exactly one trailing newline
		if !strings.HasSuffix(out, "\n") {
			t.Error("output should end with newline")
		}
		if strings.HasSuffix(out, "\n\n") {
			t.Error("output should not have trailing double newline")
		}
	})

	t.Run("modules preserve input order", func(t *testing.T) {
		// Render does not sort — it uses the order given
		mods := []*model.Module{
			{Layer: 0, Front: model.Frontmatter{ID: "base"}, Body: "L0"},
			{Layer: 1, Front: model.Frontmatter{ID: "modes/test"}, Body: "L1"},
			{Layer: 2, Front: model.Frontmatter{ID: "traits/a"}, Body: "L2"},
		}
		out, _ := Render(mods, nil)

		baseIdx := strings.Index(out, "L0")
		modeIdx := strings.Index(out, "L1")
		traitIdx := strings.Index(out, "L2")

		if baseIdx >= modeIdx || modeIdx >= traitIdx {
			t.Error("modules should preserve input order")
		}
	})

	t.Run("deduplicates unresolved variables", func(t *testing.T) {
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "base"}, Body: "{{same}}"},
			{Front: model.Frontmatter{ID: "traits/a"}, Body: "{{same}}"},
		}
		_, unresolved := Render(mods, nil)
		if len(unresolved) != 1 {
			t.Errorf("expected 1 deduplicated unresolved var, got %d: %v", len(unresolved), unresolved)
		}
	})
}

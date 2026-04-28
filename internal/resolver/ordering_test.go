package resolver

import (
	"strings"
	"testing"

	"github.com/bkuri/ppc/internal/model"
)

func TestValidateExclusiveGroups(t *testing.T) {
	t.Run("no conflicts", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{"risk"}}
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "a", Tags: []string{"risk:low"}}},
			{Front: model.Frontmatter{ID: "b", Tags: []string{"risk:low"}}},
		}

		err := ValidateExclusiveGroups(rules, mods)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("conflicting tags in same group", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{"risk"}}
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "a", Tags: []string{"risk:low"}}},
			{Front: model.Frontmatter{ID: "b", Tags: []string{"risk:high"}}},
		}

		err := ValidateExclusiveGroups(rules, mods)
		if err == nil {
			t.Fatal("expected error for conflicting tags")
		}
		if !strings.Contains(err.Error(), "conflicting") {
			t.Errorf("error = %q, want to contain 'conflicting'", err.Error())
		}
	})

	t.Run("tags in different groups", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{"risk"}}
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "a", Tags: []string{"risk:low"}}},
			{Front: model.Frontmatter{ID: "b", Tags: []string{"tone:formal"}}},
		}

		err := ValidateExclusiveGroups(rules, mods)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("no exclusive groups configured", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{}}
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "a", Tags: []string{"risk:low"}}},
			{Front: model.Frontmatter{ID: "b", Tags: []string{"risk:high"}}},
		}

		err := ValidateExclusiveGroups(rules, mods)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("empty module list", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{"risk"}}
		err := ValidateExclusiveGroups(rules, nil)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("invalid tag format", func(t *testing.T) {
		rules := &model.Rules{ExclusiveGroups: []string{"risk"}}
		mods := []*model.Module{
			{Front: model.Frontmatter{ID: "a", Tags: []string{"nocolon"}}},
		}

		err := ValidateExclusiveGroups(rules, mods)
		if err == nil {
			t.Fatal("expected error for invalid tag format")
		}
		if !strings.Contains(err.Error(), "invalid tag") {
			t.Errorf("error = %q, want to contain 'invalid tag'", err.Error())
		}
	})
}

func TestSortModules(t *testing.T) {
	t.Run("sorts by layer then priority then id", func(t *testing.T) {
		mods := []*model.Module{
			{Layer: 2, Front: model.Frontmatter{ID: "traits/b", Priority: 5}},
			{Layer: 0, Front: model.Frontmatter{ID: "base", Priority: 0}},
			{Layer: 1, Front: model.Frontmatter{ID: "modes/explore", Priority: 0}},
			{Layer: 2, Front: model.Frontmatter{ID: "traits/a", Priority: 3}},
		}

		sorted := SortModules(mods)

		if sorted[0].Front.ID != "base" {
			t.Errorf("first = %q, want 'base'", sorted[0].Front.ID)
		}
		if sorted[1].Front.ID != "modes/explore" {
			t.Errorf("second = %q, want 'modes/explore'", sorted[1].Front.ID)
		}
		if sorted[2].Front.ID != "traits/a" {
			t.Errorf("third = %q, want 'traits/a'", sorted[2].Front.ID)
		}
		if sorted[3].Front.ID != "traits/b" {
			t.Errorf("fourth = %q, want 'traits/b'", sorted[3].Front.ID)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		sorted := SortModules(nil)
		if len(sorted) != 0 {
			t.Errorf("expected empty list, got %d items", len(sorted))
		}
	})

	t.Run("single module", func(t *testing.T) {
		mods := []*model.Module{
			{Layer: 1, Front: model.Frontmatter{ID: "modes/test"}},
		}
		sorted := SortModules(mods)
		if len(sorted) != 1 || sorted[0].Front.ID != "modes/test" {
			t.Error("single module should be unchanged")
		}
	})

	t.Run("does not mutate input", func(t *testing.T) {
		mods := []*model.Module{
			{Layer: 2, Front: model.Frontmatter{ID: "traits/a"}},
			{Layer: 0, Front: model.Frontmatter{ID: "base"}},
		}
		_ = SortModules(mods)
		if mods[0].Front.ID != "traits/a" {
			t.Error("SortModules should not mutate input slice")
		}
	})
}

package resolver

import (
	"strings"
	"testing"

	"github.com/bkuri/ppc/internal/model"
)

func TestExpandRequires(t *testing.T) {
	t.Run("simple linear chain", func(t *testing.T) {
		// a requires b requires c
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"c"}}},
			"c": {Front: model.Frontmatter{ID: "c"}},
		}

		closure, fromReq, err := ExpandRequires([]string{"a"}, all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(closure) != 3 {
			t.Errorf("closure length = %d, want 3", len(closure))
		}

		// All modules should be in closure
		found := map[string]bool{}
		for _, id := range closure {
			found[id] = true
		}
		for _, id := range []string{"a", "b", "c"} {
			if !found[id] {
				t.Errorf("closure missing %q", id)
			}
		}

		// b and c should be marked as from requires (not selected)
		if !fromReq["b"] {
			t.Error("b should be marked as fromReq")
		}
		if !fromReq["c"] {
			t.Error("c should be marked as fromReq")
		}
		if fromReq["a"] {
			t.Error("a should not be marked as fromReq (it was selected)")
		}
	})

	t.Run("no requires", func(t *testing.T) {
		all := map[string]*model.Module{
			"standalone": {Front: model.Frontmatter{ID: "standalone"}},
		}

		closure, _, err := ExpandRequires([]string{"standalone"}, all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(closure) != 1 {
			t.Errorf("closure length = %d, want 1", len(closure))
		}
		if closure[0] != "standalone" {
			t.Errorf("closure[0] = %q, want %q", closure[0], "standalone")
		}
	})

	t.Run("circular dependency", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"a"}}},
		}

		_, _, err := ExpandRequires([]string{"a"}, all)
		if err == nil {
			t.Fatal("expected error for circular dependency")
		}
		if !strings.Contains(err.Error(), "circular") {
			t.Errorf("error = %q, want to contain 'circular'", err.Error())
		}
	})

	t.Run("missing dependency", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"nonexistent"}}},
		}

		_, _, err := ExpandRequires([]string{"a"}, all)
		if err == nil {
			t.Fatal("expected error for missing dependency")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("error = %q, want to contain 'not found'", err.Error())
		}
	})

	t.Run("multiple selected with shared dependency", func(t *testing.T) {
		all := map[string]*model.Module{
			"a":      {Front: model.Frontmatter{ID: "a", Requires: []string{"shared"}}},
			"b":      {Front: model.Frontmatter{ID: "b", Requires: []string{"shared"}}},
			"shared": {Front: model.Frontmatter{ID: "shared"}},
		}

		closure, _, err := ExpandRequires([]string{"a", "b"}, all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		found := map[string]bool{}
		for _, id := range closure {
			found[id] = true
		}
		for _, id := range []string{"a", "b", "shared"} {
			if !found[id] {
				t.Errorf("closure missing %q", id)
			}
		}
	})

	t.Run("diamond dependency", func(t *testing.T) {
		// a -> b -> d, a -> c -> d
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b", "c"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"d"}}},
			"c": {Front: model.Frontmatter{ID: "c", Requires: []string{"d"}}},
			"d": {Front: model.Frontmatter{ID: "d"}},
		}

		closure, _, err := ExpandRequires([]string{"a"}, all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(closure) != 4 {
			t.Errorf("closure length = %d, want 4", len(closure))
		}
	})
}

func TestDetectCycles(t *testing.T) {
	t.Run("no cycles", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"c"}}},
			"c": {Front: model.Frontmatter{ID: "c"}},
		}

		err := DetectCycles(all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("direct cycle", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"a"}}},
		}

		err := DetectCycles(all)
		if err == nil {
			t.Fatal("expected error for direct cycle")
		}
		if !strings.Contains(err.Error(), "circular") {
			t.Errorf("error = %q, want to contain 'circular'", err.Error())
		}
	})

	t.Run("three-node cycle", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"b"}}},
			"b": {Front: model.Frontmatter{ID: "b", Requires: []string{"c"}}},
			"c": {Front: model.Frontmatter{ID: "c", Requires: []string{"a"}}},
		}

		err := DetectCycles(all)
		if err == nil {
			t.Fatal("expected error for three-node cycle")
		}
	})

	t.Run("self-cycle", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"a"}}},
		}

		err := DetectCycles(all)
		if err == nil {
			t.Fatal("expected error for self-cycle")
		}
	})

	t.Run("missing dependency is skipped", func(t *testing.T) {
		all := map[string]*model.Module{
			"a": {Front: model.Frontmatter{ID: "a", Requires: []string{"missing"}}},
		}

		err := DetectCycles(all)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("empty module set", func(t *testing.T) {
		err := DetectCycles(map[string]*model.Module{})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

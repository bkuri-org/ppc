package loader

import (
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	t.Run("valid frontmatter with fields", func(t *testing.T) {
		raw := []byte("---\nid: test\ndesc: A test module\ntags:\n  - risk:low\n---\nBody content here.\n")
		fm, body, has, err := ParseFrontmatter(raw)
		if err.Msg != "" {
			t.Fatalf("unexpected error: %s", err.Msg)
		}
		if !has {
			t.Fatal("expected has=true")
		}
		if fm.ID != "test" {
			t.Errorf("ID = %q, want %q", fm.ID, "test")
		}
		if fm.Desc != "A test module" {
			t.Errorf("Desc = %q, want %q", fm.Desc, "A test module")
		}
		if len(fm.Tags) != 1 || fm.Tags[0] != "risk:low" {
			t.Errorf("Tags = %v, want [risk:low]", fm.Tags)
		}
		if body != "Body content here." {
			t.Errorf("Body = %q, want %q", body, "Body content here.")
		}
	})

	t.Run("empty YAML frontmatter", func(t *testing.T) {
		raw := []byte("---\n\n---\nSome body.\n")
		fm, body, has, err := ParseFrontmatter(raw)
		if err.Msg != "" {
			t.Fatalf("unexpected error: %s", err.Msg)
		}
		if !has {
			t.Fatal("expected has=true")
		}
		if fm.ID != "" {
			t.Errorf("ID = %q, want empty", fm.ID)
		}
		if body != "Some body." {
			t.Errorf("Body = %q, want %q", body, "Some body.")
		}
	})

	t.Run("no frontmatter delimiters", func(t *testing.T) {
		raw := []byte("Just plain content\nwith no frontmatter.\n")
		fm, body, has, err := ParseFrontmatter(raw)
		if err.Msg != "" {
			t.Fatalf("unexpected error: %s", err.Msg)
		}
		if has {
			t.Fatal("expected has=false")
		}
		if fm.ID != "" {
			t.Errorf("ID = %q, want empty", fm.ID)
		}
		if !strings.Contains(body, "Just plain content") {
			t.Errorf("Body = %q, want to contain 'Just plain content'", body)
		}
	})

	t.Run("missing closing delimiter", func(t *testing.T) {
		raw := []byte("---\nid: test\nNo closing delimiter\n")
		_, _, has, err := ParseFrontmatter(raw)
		if err.Msg == "" {
			t.Fatal("expected error for missing closing delimiter")
		}
		if has {
			t.Fatal("expected has=false")
		}
	})

	t.Run("invalid YAML in frontmatter", func(t *testing.T) {
		raw := []byte("---\nid: [invalid yaml\n---\nBody\n")
		_, _, _, err := ParseFrontmatter(raw)
		if err.Msg == "" {
			t.Fatal("expected error for invalid YAML")
		}
		if !strings.Contains(err.Msg, "invalid YAML frontmatter") {
			t.Errorf("error = %q, want to contain 'invalid YAML frontmatter'", err.Msg)
		}
	})

	t.Run("frontmatter with requires", func(t *testing.T) {
		raw := []byte("---\nid: traits/deep\nrequires:\n  - base\n---\nTrait body.\n")
		fm, _, has, err := ParseFrontmatter(raw)
		if err.Msg != "" {
			t.Fatalf("unexpected error: %s", err.Msg)
		}
		if !has {
			t.Fatal("expected has=true")
		}
		if len(fm.Requires) != 1 || fm.Requires[0] != "base" {
			t.Errorf("Requires = %v, want [base]", fm.Requires)
		}
	})

	t.Run("frontmatter with priority", func(t *testing.T) {
		raw := []byte("---\nid: policies/rule\npriority: 10\n---\nRule body.\n")
		fm, _, has, err := ParseFrontmatter(raw)
		if err.Msg != "" {
			t.Fatalf("unexpected error: %s", err.Msg)
		}
		if !has {
			t.Fatal("expected has=true")
		}
		if fm.Priority != 10 {
			t.Errorf("Priority = %d, want 10", fm.Priority)
		}
	})
}

func TestLoadModules(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		modByID, err := LoadModules("testdata/valid")
		if err != nil {
			t.Fatalf("LoadModules failed: %s", err)
		}
		if len(modByID) != 2 {
			t.Errorf("got %d modules, want 2", len(modByID))
		}
		if _, ok := modByID["base"]; !ok {
			t.Error("missing module 'base'")
		}
		if _, ok := modByID["traits/deep"]; !ok {
			t.Error("missing module 'traits/deep'")
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		// ListMarkdownFiles silently returns empty for nonexistent dirs,
		// so LoadModules succeeds with an empty map
		modByID, err := LoadModules("testdata/nonexistent_dir")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(modByID) != 0 {
			t.Errorf("expected empty map, got %d modules", len(modByID))
		}
	})

	t.Run("directory with non-.md files only", func(t *testing.T) {
		_, err := LoadModules("testdata/skip_nonmd")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("missing frontmatter", func(t *testing.T) {
		_, err := LoadModules("testdata/missing_frontmatter")
		if err == nil {
			t.Fatal("expected error for missing frontmatter")
		}
		if !strings.Contains(err.Error(), "missing frontmatter") {
			t.Errorf("error = %q, want to contain 'missing frontmatter'", err.Error())
		}
	})
}

func TestListMarkdownFiles(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		files := ListMarkdownFiles("testdata/valid")
		if len(files) == 0 {
			t.Fatal("expected at least one .md file")
		}
		for _, f := range files {
			if !strings.HasSuffix(f, ".md") {
				t.Errorf("found non-.md file: %s", f)
			}
		}
	})

	t.Run("nonexistent directory returns empty", func(t *testing.T) {
		files := ListMarkdownFiles("testdata/nonexistent_dir_xyz")
		if len(files) != 0 {
			t.Errorf("expected empty slice for nonexistent dir, got %d files", len(files))
		}
	})

	t.Run("skips non-.md files", func(t *testing.T) {
		files := ListMarkdownFiles("testdata/skip_nonmd")
		if len(files) != 1 {
			t.Errorf("expected 1 .md file, got %d", len(files))
		}
		if !strings.HasSuffix(files[0], "base.md") {
			t.Errorf("expected base.md, got %s", files[0])
		}
	})
}

func TestLoadRules(t *testing.T) {
	t.Run("valid rules file", func(t *testing.T) {
		rules, err := LoadRules("testdata/valid")
		if err != nil {
			t.Fatalf("LoadRules failed: %s", err)
		}
		if len(rules.ExclusiveGroups) != 1 || rules.ExclusiveGroups[0] != "risk" {
			t.Errorf("ExclusiveGroups = %v, want [risk]", rules.ExclusiveGroups)
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := LoadRules("testdata/nonexistent_dir")
		if err == nil {
			t.Fatal("expected error for nonexistent directory")
		}
		if !strings.Contains(err.Error(), "missing rules file") {
			t.Errorf("error = %q, want to contain 'missing rules file'", err.Error())
		}
	})

	t.Run("invalid YAML in rules file", func(t *testing.T) {
		_, err := LoadRules("testdata/invalid_rules")
		if err == nil {
			t.Fatal("expected error for invalid YAML")
		}
		if !strings.Contains(err.Error(), "invalid rules.yml") {
			t.Errorf("error = %q, want to contain 'invalid rules.yml'", err.Error())
		}
	})
}

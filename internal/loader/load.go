// Package loader handles loading modules and rules from the filesystem.
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	errtypes "github.com/bkuri/ppc/internal/error"
	"github.com/bkuri/ppc/internal/model"
	"gopkg.in/yaml.v3"
)

// ListMarkdownFiles finds all .md files recursively
func ListMarkdownFiles(root string) []string {
	var out []string
	_ = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// LoadModules loads all module files from promptsDir
func LoadModules(promptsDir string) (map[string]*model.Module, error) {
	files := ListMarkdownFiles(promptsDir)
	modByID := map[string]*model.Module{}

	for _, p := range files {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, errtypes.New(p, "", fmt.Sprintf("failed to read: %v", err))
		}
		fm, body, has, parseErr := ParseFrontmatter(raw)
		if parseErr.Msg != "" {
			return nil, parseErr
		}
		if !has {
			return nil, errtypes.New(p, "", "missing frontmatter (v0.1 requires YAML frontmatter with id)")
		}
		if strings.TrimSpace(fm.ID) == "" {
			return nil, errtypes.New(p, "", "frontmatter missing required field: id")
		}
		m := &model.Module{
			Path:  p,
			Layer: model.LayerIndexFromPath(filepath.ToSlash(p)),
			Front: fm,
			Body:  body,
		}
		if _, exists := modByID[fm.ID]; exists {
			return nil, errtypes.New(p, fm.ID, fmt.Sprintf("duplicate module id %q", fm.ID))
		}
		modByID[fm.ID] = m
	}
	return modByID, nil
}

// LoadRules loads the rules.yml file from promptsDir
func LoadRules(promptsDir string) (*model.Rules, error) {
	p := filepath.Join(promptsDir, "rules.yml")
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, errtypes.New(p, "", fmt.Sprintf("missing rules file: %v", err))
	}
	var r model.Rules
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, errtypes.New(p, "", fmt.Sprintf("invalid rules.yml: %v", err))
	}
	return &r, nil
}

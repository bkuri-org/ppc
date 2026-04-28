// Package model defines core data structures
package model

import (
	"path/filepath"
	"strings"
)

// LintContentPattern defines a content-level lint rule
type LintContentPattern struct {
	Match  string   `yaml:"match"`
	Reason string   `yaml:"reason"`
	Paths  []string `yaml:"paths,omitempty"`
}

// LintScope defines path-scoped lint overrides
type LintScope struct {
	Paths          []string             `yaml:"paths"`
	MaxModuleWords *int                 `yaml:"max_module_words,omitempty"`
	ForbidEmpty    *bool                `yaml:"forbid_empty_body,omitempty"`
	RequireFields  []string             `yaml:"require_fields,omitempty"`
	ForbidTags     []string             `yaml:"forbid_tags,omitempty"`
	ContentPattern []LintContentPattern `yaml:"forbid_content_patterns,omitempty"`
}

// LintConfig defines persistent lint configuration
type LintConfig struct {
	MaxWords              int                  `yaml:"max_words"`
	MaxLines              int                  `yaml:"max_lines"`
	MaxModules            int                  `yaml:"max_modules"`
	MaxModuleWords        int                  `yaml:"max_module_words"`
	MaxDepth              int                  `yaml:"max_depth"`
	RequireTags           []string             `yaml:"require_tags"`
	ForbidTags            []string             `yaml:"forbid_tags"`
	RequireFields         []string             `yaml:"require_fields"`
	ForbidEmptyBody       bool                 `yaml:"forbid_empty_body"`
	ForbidContentPatterns []LintContentPattern `yaml:"forbid_content_patterns"`
	Scopes                []LintScope          `yaml:"scopes"`
}

// Rules defines validation rules for modules
type Rules struct {
	ExclusiveGroups []string   `yaml:"exclusive_groups"`
	Lint            LintConfig `yaml:"lint"`
}

// Frontmatter represents the YAML frontmatter of a module
type Frontmatter struct {
	ID       string   `yaml:"id"`
	Desc     string   `yaml:"desc"`
	Priority int      `yaml:"priority"`
	Tags     []string `yaml:"tags"`
	Requires []string `yaml:"requires"`
}

// Module represents a compiled module with metadata
type Module struct {
	Path     string
	Layer    int
	Front    Frontmatter
	Body     string
	FromReq  bool
	Selected bool
}

// LayerOrder defines the canonical layer precedence
var LayerOrder = []string{"base", "modes", "traits", "policies", "contracts", "guardrails"}

func LayerName(index int) string {
	if index < 0 || index >= len(LayerOrder) {
		return "unknown"
	}
	return LayerOrder[index]
}

// LayerIndexFromPath returns the layer index for a module path
func LayerIndexFromPath(p string) int {
	parts := strings.Split(filepath.ToSlash(p), "/")
	for i, s := range LayerOrder {
		for _, part := range parts {
			if part == s {
				return i
			}
		}
	}
	return 0
}

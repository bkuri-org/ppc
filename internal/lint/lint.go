package lint

import (
	"fmt"
	"strings"

	"github.com/bkuri/ppc/internal/loader"
	"github.com/bkuri/ppc/internal/model"
)

type Config struct {
	MaxWords        int
	MaxLines        int
	MaxModules      int
	MaxModuleWords  int
	MaxDepth        int
	RequireTags     []string
	ForbidTags      []string
	RequireFields   []string
	ForbidEmptyBody bool
}

type Violation struct {
	Level   string `json:"level"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Module  string `json:"module,omitempty"`
}

type Result struct {
	Violations []Violation    `json:"violations"`
	Stats      map[string]int `json:"stats"`
}

func Run(promptsDir string, cfg Config) (*Result, error) {
	modByID, errIf := loader.LoadModules(promptsDir)
	if errIf != nil {
		return nil, errIf.(error)
	}

	result := &Result{
		Violations: []Violation{},
		Stats:      make(map[string]int),
	}

	result.Stats["module_count"] = len(modByID)

	totalWords := 0
	totalLines := 0
	maxModuleWords := 0

	for _, m := range modByID {
		words := countWords(m.Body)
		lines := countLines(m.Body)
		totalWords += words
		totalLines += lines
		if words > maxModuleWords {
			maxModuleWords = words
		}
	}

	result.Stats["word_count"] = totalWords
	result.Stats["line_count"] = totalLines
	result.Stats["max_module_words"] = maxModuleWords

	if cfg.MaxWords > 0 && totalWords > cfg.MaxWords {
		pct := percentOver(totalWords, cfg.MaxWords)
		result.Violations = append(result.Violations, Violation{
			Level:   "WARN",
			Rule:    "max_words",
			Message: fmt.Sprintf("word count (%d) exceeds threshold (%d) by %d%%", totalWords, cfg.MaxWords, pct),
			Module:  "",
		})
	}

	if cfg.MaxLines > 0 && totalLines > cfg.MaxLines {
		pct := percentOver(totalLines, cfg.MaxLines)
		result.Violations = append(result.Violations, Violation{
			Level:   "WARN",
			Rule:    "max_lines",
			Message: fmt.Sprintf("line count (%d) exceeds threshold (%d) by %d%%", totalLines, cfg.MaxLines, pct),
			Module:  "",
		})
	}

	if cfg.MaxModules > 0 && len(modByID) > cfg.MaxModules {
		pct := percentOver(len(modByID), cfg.MaxModules)
		result.Violations = append(result.Violations, Violation{
			Level:   "WARN",
			Rule:    "max_modules",
			Message: fmt.Sprintf("module count (%d) exceeds threshold (%d) by %d%%", len(modByID), cfg.MaxModules, pct),
			Module:  "",
		})
	}

	for id, m := range modByID {
		if cfg.MaxModuleWords > 0 {
			words := countWords(m.Body)
			if words > cfg.MaxModuleWords {
				pct := percentOver(words, cfg.MaxModuleWords)
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "max_module_words",
					Message: fmt.Sprintf("word count (%d) exceeds threshold (%d) by %d%%", words, cfg.MaxModuleWords, pct),
					Module:  id,
				})
			}
		}

		if cfg.ForbidEmptyBody && strings.TrimSpace(m.Body) == "" {
			result.Violations = append(result.Violations, Violation{
				Level:   "WARN",
				Rule:    "forbid_empty_body",
				Message: "module has empty body",
				Module:  id,
			})
		}

		for _, field := range cfg.RequireFields {
			if !hasField(m.Front, field) {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "require_fields",
					Message: "missing required field '" + field + "'",
					Module:  id,
				})
			}
		}

		for _, ft := range cfg.ForbidTags {
			for _, t := range m.Front.Tags {
				if t == ft {
					result.Violations = append(result.Violations, Violation{
						Level:   "WARN",
						Rule:    "forbid_tags",
						Message: "module has forbidden tag '" + ft + "'",
						Module:  id,
					})
				}
			}
		}
	}

	if len(cfg.RequireTags) > 0 {
		allTags := []string{}
		for _, m := range modByID {
			allTags = append(allTags, m.Front.Tags...)
		}
		for _, pattern := range cfg.RequireTags {
			if !tagPatternMatches(pattern, allTags) {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "require_tags",
					Message: "no module has required tag pattern '" + pattern + "'",
					Module:  "",
				})
			}
		}
	}

	if cfg.MaxDepth > 0 {
		for id := range modByID {
			depth, chain := calculateModuleDepth(modByID, id)
			if depth > cfg.MaxDepth {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "max_depth",
					Message: formatDepthMessage(depth, cfg.MaxDepth, chain),
					Module:  id,
				})
			}
		}
	}

	return result, nil
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func percentOver(actual, threshold int) int {
	if threshold == 0 {
		return 0
	}
	return ((actual - threshold) * 100) / threshold
}

func hasField(fm model.Frontmatter, field string) bool {
	switch field {
	case "id":
		return fm.ID != ""
	case "desc":
		return fm.Desc != ""
	case "priority":
		return fm.Priority != 0
	case "tags":
		return len(fm.Tags) > 0
	case "requires":
		return len(fm.Requires) > 0
	default:
		return false
	}
}

func calculateModuleDepth(modByID map[string]*model.Module, startID string) (int, []string) {
	depthMemo := make(map[string]int)
	chainMemo := make(map[string][]string)

	var getDepth func(id string, path []string) (int, []string)
	getDepth = func(id string, path []string) (int, []string) {
		for _, p := range path {
			if p == id {
				return len(path), append(path, id)
			}
		}

		if d, exists := depthMemo[id]; exists {
			return d, chainMemo[id]
		}

		m, ok := modByID[id]
		if !ok {
			return len(path), path
		}

		if len(m.Front.Requires) == 0 {
			depthMemo[id] = 1
			chainMemo[id] = []string{id}
			return 1, []string{id}
		}

		maxDepth := 0
		var maxChain []string

		for _, req := range m.Front.Requires {
			if _, ok := modByID[req]; !ok {
				continue
			}
			d, chain := getDepth(req, append(path, id))
			if d > maxDepth {
				maxDepth = d
				maxChain = chain
			}
		}

		depth := maxDepth + 1
		fullChain := append([]string{id}, maxChain...)

		depthMemo[id] = depth
		chainMemo[id] = fullChain

		return depth, fullChain
	}

	return getDepth(startID, []string{})
}

func formatDepthMessage(depth, threshold int, chain []string) string {
	return fmt.Sprintf("dependency depth (%d) exceeds threshold (%d): chain = %s", depth, threshold, strings.Join(chain, " -> "))
}

func tagPatternMatches(pattern string, tags []string) bool {
	if strings.HasSuffix(pattern, ":*") {
		group := strings.TrimSuffix(pattern, ":*")
		for _, t := range tags {
			if strings.HasPrefix(t, group+":") {
				return true
			}
		}
		return false
	}
	for _, t := range tags {
		if t == pattern {
			return true
		}
	}
	return false
}

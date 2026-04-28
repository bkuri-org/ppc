package lint

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bkuri/ppc/internal/loader"
	"github.com/bkuri/ppc/internal/model"
)

type Config struct {
	MaxWords               int
	MaxLines               int
	MaxModules             int
	MaxModuleWords         int
	MaxDepth               int
	RequireTags            []string
	ForbidTags             []string
	RequireFields          []string
	ForbidEmptyBody        bool
	ForbidContentPatterns  []ContentPattern
}

type ContentPattern struct {
	Match  string
	Reason string
	Paths  []string
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

func MergeConfig(file model.LintConfig, cli Config) Config {
	merged := Config{
		MaxWords:      coalesceInt(file.MaxWords, cli.MaxWords),
		MaxLines:      coalesceInt(file.MaxLines, cli.MaxLines),
		MaxModules:    coalesceInt(file.MaxModules, cli.MaxModules),
		MaxModuleWords: coalesceInt(file.MaxModuleWords, cli.MaxModuleWords),
		MaxDepth:      coalesceInt(file.MaxDepth, cli.MaxDepth),
		ForbidEmptyBody: coalesceBool(file.ForbidEmptyBody, cli.ForbidEmptyBody),
	}

	if len(cli.RequireTags) > 0 {
		merged.RequireTags = cli.RequireTags
	} else if len(file.RequireTags) > 0 {
		merged.RequireTags = file.RequireTags
	}

	if len(cli.ForbidTags) > 0 {
		merged.ForbidTags = cli.ForbidTags
	} else if len(file.ForbidTags) > 0 {
		merged.ForbidTags = file.ForbidTags
	}

	if len(cli.RequireFields) > 0 {
		merged.RequireFields = cli.RequireFields
	} else if len(file.RequireFields) > 0 {
		merged.RequireFields = file.RequireFields
	}

	if len(cli.ForbidContentPatterns) > 0 {
		merged.ForbidContentPatterns = cli.ForbidContentPatterns
	} else if len(file.ForbidContentPatterns) > 0 {
		for _, p := range file.ForbidContentPatterns {
			merged.ForbidContentPatterns = append(merged.ForbidContentPatterns, ContentPattern{
				Match:  p.Match,
				Reason: p.Reason,
				Paths:  p.Paths,
			})
		}
	}

	return merged
}

func coalesceInt(file, cli int) int {
	if cli != 0 {
		return cli
	}
	return file
}

func coalesceBool(file, cli bool) bool {
	if cli {
		return true
	}
	return file
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
		})
	}

	if cfg.MaxLines > 0 && totalLines > cfg.MaxLines {
		pct := percentOver(totalLines, cfg.MaxLines)
		result.Violations = append(result.Violations, Violation{
			Level:   "WARN",
			Rule:    "max_lines",
			Message: fmt.Sprintf("line count (%d) exceeds threshold (%d) by %d%%", totalLines, cfg.MaxLines, pct),
		})
	}

	if cfg.MaxModules > 0 && len(modByID) > cfg.MaxModules {
		pct := percentOver(len(modByID), cfg.MaxModules)
		result.Violations = append(result.Violations, Violation{
			Level:   "WARN",
			Rule:    "max_modules",
			Message: fmt.Sprintf("module count (%d) exceeds threshold (%d) by %d%%", len(modByID), cfg.MaxModules, pct),
		})
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

	sortedIDs := make([]string, 0, len(modByID))
	for id := range modByID {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Strings(sortedIDs)

	for _, id := range sortedIDs {
		m := modByID[id]
		scope := resolveScope(m.Path, cfg)

		if scope.MaxModuleWords > 0 {
			words := countWords(m.Body)
			if words > scope.MaxModuleWords {
				pct := percentOver(words, scope.MaxModuleWords)
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "max_module_words",
					Message: fmt.Sprintf("word count (%d) exceeds threshold (%d) by %d%%", words, scope.MaxModuleWords, pct),
					Module:  id,
				})
			}
		}

		if scope.ForbidEmptyBody && strings.TrimSpace(m.Body) == "" {
			result.Violations = append(result.Violations, Violation{
				Level:   "WARN",
				Rule:    "forbid_empty_body",
				Message: "module has empty body",
				Module:  id,
			})
		}

		for _, field := range scope.RequireFields {
			if !hasField(m.Front, field) {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "require_fields",
					Message: "missing required field '" + field + "'",
					Module:  id,
				})
			}
		}

		for _, ft := range scope.ForbidTags {
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

		for _, cp := range cfg.ForbidContentPatterns {
			if len(cp.Paths) > 0 && !matchPaths(m.Path, cp.Paths) {
				continue
			}
			re, err := regexp.Compile(cp.Match)
			if err != nil {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "forbid_content",
					Message: fmt.Sprintf("invalid pattern %q: %v", cp.Match, err),
					Module:  id,
				})
				continue
			}
			if re.MatchString(m.Body) {
				result.Violations = append(result.Violations, Violation{
					Level:   "WARN",
					Rule:    "forbid_content",
					Message: cp.Reason,
					Module:  id,
				})
			}
		}
	}

	return result, nil
}

type resolvedScope struct {
	MaxModuleWords int
	ForbidEmptyBody bool
	RequireFields  []string
	ForbidTags     []string
}

func resolveScope(modPath string, cfg Config) resolvedScope {
	scoped := resolvedScope{
		MaxModuleWords: cfg.MaxModuleWords,
		ForbidEmptyBody: cfg.ForbidEmptyBody,
		RequireFields:  cfg.RequireFields,
		ForbidTags:     cfg.ForbidTags,
	}

	return scoped
}

func matchPaths(modPath string, patterns []string) bool {
	for _, pat := range patterns {
		if matchGlob(modPath, pat) {
			return true
		}
	}
	return false
}

func matchGlob(path, pattern string) bool {
	if pattern == "" {
		return false
	}
	if pattern == "**" {
		return true
	}

	relPath := filepath.ToSlash(path)
	relPat := filepath.ToSlash(pattern)

	parts := strings.Split(relPat, "/")
	pathParts := strings.Split(relPath, "/")

	return globMatch(pathParts, parts)
}

func globMatch(pathParts, patParts []string) bool {
	if len(patParts) == 0 {
		return len(pathParts) == 0
	}

	seg := patParts[0]
	rest := patParts[1:]

	if seg == "**" {
		if len(rest) == 0 {
			return true
		}
		for i := 0; i <= len(pathParts); i++ {
			if globMatch(pathParts[i:], rest) {
				return true
			}
		}
		return false
	}

	if len(pathParts) == 0 {
		return false
	}

	matched, _ := filepath.Match(seg, pathParts[0])
	if !matched {
		return false
	}

	return globMatch(pathParts[1:], rest)
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

			if d, exists := depthMemo[id]; exists {
				return d, chainMemo[id]
			}
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

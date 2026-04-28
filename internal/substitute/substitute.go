package substitute

import (
	"fmt"
	"regexp"
	"strings"
)

type Vars map[string]any

var varPattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func Substitute(content string, vars Vars) (string, []string) {
	var unresolved []string
	result := varPattern.ReplaceAllStringFunc(content, func(match string) string {
		path := strings.TrimSpace(varPattern.FindStringSubmatch(match)[1])
		val, ok := ResolvePath(vars, path)
		if !ok {
			unresolved = append(unresolved, path)
			return match
		}
		return formatValue(val)
	})
	return result, unresolved
}

func ResolvePath(vars Vars, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = vars

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, false
			}
		case Vars:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, false
			}
		default:
			return nil, false
		}
	}

	return current, true
}

func formatValue(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

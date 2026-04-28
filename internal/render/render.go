package render

import (
	"strings"

	"github.com/bkuri/ppc/internal/model"
	"github.com/bkuri/ppc/internal/substitute"
)

func Render(mods []*model.Module, vars substitute.Vars) (string, []string) {
	var b strings.Builder
	seen := map[string]bool{}
	var unresolved []string
	for i, m := range mods {
		if i > 0 {
			b.WriteString("\n\n")
		}
		body, unres := substitute.Substitute(m.Body, vars)
		for _, u := range unres {
			if !seen[u] {
				seen[u] = true
				unresolved = append(unresolved, u)
			}
		}
		b.WriteString(strings.TrimRight(body, "\n"))
	}
	return strings.TrimRight(b.String(), "\n") + "\n", unresolved
}

// Package resolver handles module dependency resolution and ordering.
package resolver

import (
	"fmt"
	"sort"
	"strings"

	errtypes "github.com/bkuri/ppc/internal/error"
	"github.com/bkuri/ppc/internal/model"
)

// ExpandRequires performs transitive closure of requires dependencies
// Returns (closureIDs, fromReq map, error)
func ExpandRequires(selectedIDs []string, all map[string]*model.Module) ([]string, map[string]bool, error) {
	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)

	state := map[string]int{}
	stack := []string{}
	pos := map[string]int{}
	out := []string{}
	inOut := map[string]bool{}
	fromReq := map[string]bool{}

	var dfs func(id string, rootSelected bool) error
	dfs = func(id string, rootSelected bool) error {
		m, ok := all[id]
		if !ok {
			return errtypes.New("", id, fmt.Sprintf("required module not found: %s", id))
		}

		switch state[id] {
		case done:
			if rootSelected {
				m.Selected = true
				fromReq[id] = false
			}
			return nil
		case visiting:
			i := pos[id]
			cycle := append(append([]string{}, stack[i:]...), id)
			return errtypes.New("", id, fmt.Sprintf("circular requires: %s", strings.Join(cycle, " -> ")))
		}

		state[id] = visiting
		pos[id] = len(stack)
		stack = append(stack, id)

		reqs := append([]string{}, m.Front.Requires...)
		sort.Strings(reqs)
		for _, r := range reqs {
			if err := dfs(r, false); err != nil {
				return err
			}
			if !Contains(selectedIDs, r) {
				fromReq[r] = true
			}
		}

		stack = stack[:len(stack)-1]
		delete(pos, id)
		state[id] = done

		if !inOut[id] {
			inOut[id] = true
			out = append(out, id)
		}
		if rootSelected {
			m.Selected = true
			fromReq[id] = false
		}
		return nil
	}

	ids := append([]string{}, selectedIDs...)
	sort.Strings(ids)
	for _, id := range ids {
		if err := dfs(id, true); err != nil {
			return nil, nil, err
		}
	}

	return out, fromReq, nil
}

func DetectCycles(all map[string]*model.Module) error {
	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)

	state := map[string]int{}
	stack := []string{}
	pos := map[string]int{}

	var visit func(id string) error
	visit = func(id string) error {
		switch state[id] {
		case done:
			return nil
		case visiting:
			i := pos[id]
			cycle := append(append([]string{}, stack[i:]...), id)
			return fmt.Errorf("circular requires: %s", strings.Join(cycle, " -> "))
		}

		state[id] = visiting
		pos[id] = len(stack)
		stack = append(stack, id)

		reqs := append([]string{}, all[id].Front.Requires...)
		sort.Strings(reqs)
		for _, r := range reqs {
			if _, ok := all[r]; !ok {
				continue
			}
			if err := visit(r); err != nil {
				return err
			}
		}

		stack = stack[:len(stack)-1]
		delete(pos, id)
		state[id] = done
		return nil
	}

	ids := make([]string, 0, len(all))
	for id := range all {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}

	return nil
}

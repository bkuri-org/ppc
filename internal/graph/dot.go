package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bkuri/ppc/internal/model"
)

// Edge represents a directed edge in the dependency graph
type Edge struct {
	From, To string
}

// BuildDOT generates Graphviz DOT representation of module dependency graph.
// Input: all modules, rules (for exclusive groups), reachability map
// Output: deterministic DOT string (no timestamps, randomness, or churn)
//
// Determinism rules:
// - modules sorted by (layer, id)
// - edges sorted lexicographically (source, then target)
// - attribute ordering consistent (shape, style, color)
// - subgraph names stable (cluster_0_base, cluster_1_modes, etc.)
func BuildDOT(modByID map[string]*model.Module, _ *model.Rules, reachable map[string]bool) string {
	sortedIDs := sortedModuleIDs(modByID)
	byLayer := layerSubgraphs(sortedIDs, modByID)
	edges := collectEdges(sortedIDs, modByID)

	var buf strings.Builder
	buf.WriteString("digraph ppc {\n")
	buf.WriteString("  rankdir=LR;\n\n")

	for layerIdx := 0; layerIdx < len(model.LayerOrder); layerIdx++ {
		if ids, ok := byLayer[layerIdx]; ok && len(ids) > 0 {
			clusterName := fmt.Sprintf("cluster_%d_%s", layerIdx, model.LayerName(layerIdx))
			layerLabel := model.LayerName(layerIdx)
			buf.WriteString(fmt.Sprintf("  subgraph %s {\n", clusterName))
			buf.WriteString(fmt.Sprintf("    label=\"%s\";\n", layerLabel))

			for _, id := range ids {
				buf.WriteString(fmt.Sprintf("    \"%s\";\n", id))
			}
			buf.WriteString("  }\n\n")
		}
	}

	for _, edge := range edges {
		buf.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", edge.From, edge.To))
	}

	buf.WriteString("\n")
	for _, id := range sortedIDs {
		if !reachable[id] {
			buf.WriteString(fmt.Sprintf("  \"%s\" [style=\"dashed\", color=\"red\"];\n", id))
		}
	}

	buf.WriteString("\n")
	for _, id := range sortedIDs {
		if isEntrypoint(id) {
			buf.WriteString(fmt.Sprintf("  \"%s\" [shape=\"box\", style=\"bold\"];\n", id))
		}
	}

	buf.WriteString("}\n")
	return buf.String()
}

func sortedModuleIDs(modByID map[string]*model.Module) []string {
	ids := make([]string, 0, len(modByID))
	for id := range modByID {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		layerI := modByID[ids[i]].Layer
		layerJ := modByID[ids[j]].Layer
		if layerI != layerJ {
			return layerI < layerJ
		}
		return ids[i] < ids[j]
	})
	return ids
}

func layerSubgraphs(sortedIDs []string, modByID map[string]*model.Module) map[int][]string {
	byLayer := make(map[int][]string)
	for _, id := range sortedIDs {
		layer := modByID[id].Layer
		byLayer[layer] = append(byLayer[layer], id)
	}
	return byLayer
}

func collectEdges(sortedIDs []string, modByID map[string]*model.Module) []Edge {
	var edges []Edge
	seen := make(map[string]bool)

	for _, id := range sortedIDs {
		m := modByID[id]
		for _, req := range m.Front.Requires {
			edgeKey := id + "->" + req
			if !seen[edgeKey] {
				edges = append(edges, Edge{From: id, To: req})
				seen[edgeKey] = true
			}
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		return edges[i].To < edges[j].To
	})

	return edges
}

func isEntrypoint(id string) bool {
	return id == "base" ||
		strings.HasPrefix(id, "modes/") ||
		strings.HasPrefix(id, "contracts/")
}

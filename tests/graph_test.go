package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bkuri/ppc/internal/graph"
	"github.com/bkuri/ppc/internal/loader"
	"github.com/bkuri/ppc/internal/model"
)

var promptsDir = filepath.Join("..", "prompts")

func TestGraphDeterministic(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	t.Logf("TestGraphDeterministic: loaded %d modules", len(modByID))
	reachable := computeReachable(modByID)

	dot1 := graph.BuildDOT(modByID, rules, reachable)
	dot2 := graph.BuildDOT(modByID, rules, reachable)

	if dot1 != dot2 {
		t.Fatal("graph output not deterministic")
	}
}

func TestGraphValidDOT(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	t.Logf("TestGraphValidDOT: loaded %d modules", len(modByID))
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)

	if !strings.HasPrefix(dot, "digraph ppc {") {
		t.Error("DOT does not start with 'digraph ppc {'")
	}
	if !strings.HasSuffix(dot, "}\n") {
		t.Error("DOT does not end with '}'")
	}
}

func TestGraphContainsAllModules(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)

	for id := range modByID {
		if !strings.Contains(dot, "\""+id+"\"") {
			t.Errorf("module %q not in graph", id)
		}
	}
}

func TestGraphUnreachableStyled(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)

	for id := range modByID {
		if !reachable[id] {
			pattern := `"` + id + `".*style="dashed".*color="red"`
			if !regexp.MustCompile(pattern).MatchString(dot) {
				t.Errorf("unreachable module %q not styled red+dashed", id)
			}
		}
	}
}

func TestGraphEntrypointsStyled(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)

	if !strings.Contains(dot, `shape="box"`) {
		t.Error("missing shape=box attribute")
	}
	if !strings.Contains(dot, `style="bold"`) {
		t.Error("missing style=bold attribute")
	}
}

func TestGraphGoldenSnapshot(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)
	expected, err := os.ReadFile("testdata/doctor_graph.dot")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	t.Logf("Generated DOT length: %d", len(dot))

	if dot != string(expected) {
		t.Fatalf("graph snapshot mismatch.\nGot:\n%s\n\nExpected:\n%s", dot, string(expected))
	}
}

func TestGraphHasAllLayerClusters(t *testing.T) {
	modByID, _ := loader.LoadModules(promptsDir)
	rules, _ := loader.LoadRules(promptsDir)
	reachable := computeReachable(modByID)

	dot := graph.BuildDOT(modByID, rules, reachable)

	if !strings.Contains(dot, "subgraph") || !strings.Contains(dot, "cluster") {
		t.Error("missing cluster subgraphs")
	}
}

func computeReachable(modByID map[string]*model.Module) map[string]bool {
	entry := map[string]bool{"base": true}
	for id := range modByID {
		if strings.HasPrefix(id, "modes/") || strings.HasPrefix(id, "contracts/") {
			entry[id] = true
		}
	}

	reachable := map[string]bool{}
	var mark func(id string)
	mark = func(id string) {
		if reachable[id] {
			return
		}
		reachable[id] = true
		for _, r := range modByID[id].Front.Requires {
			if _, ok := modByID[r]; ok {
				mark(r)
			}
		}
	}

	entryIDs := make([]string, 0, len(entry))
	for id := range entry {
		entryIDs = append(entryIDs, id)
	}

	for _, id := range entryIDs {
		if _, ok := modByID[id]; ok {
			mark(id)
		}
	}

	return reachable
}

package compile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bkuri/ppc/internal/loader"
	"github.com/bkuri/ppc/internal/model"
	"github.com/bkuri/ppc/internal/render"
	"github.com/bkuri/ppc/internal/resolver"
	"github.com/bkuri/ppc/internal/substitute"
	"gopkg.in/yaml.v3"
)

func Compile(opts CompileOptions) (string, CompileMeta, error) {
	modByID, err := loader.LoadModules(opts.PromptsDir)
	if err != nil {
		return "", CompileMeta{}, err
	}

	rules, err := loader.LoadRules(opts.PromptsDir)
	if err != nil {
		return "", CompileMeta{}, err
	}

	vars := substitute.Vars{}
	if opts.VarsFile != "" {
		vars, err = loadVarsFile(opts.VarsFile)
		if err != nil {
			return "", CompileMeta{}, err
		}
	}
	for k, v := range opts.Vars {
		vars[k] = v
	}

	selectedIDs := buildSelectedIDs(opts)

	closureIDs, fromReq, err := resolver.ExpandRequires(selectedIDs, modByID)
	if err != nil {
		return "", CompileMeta{}, err
	}

	mods, order := buildModuleList(closureIDs, fromReq, selectedIDs, modByID)

	if err := resolver.ValidateExclusiveGroups(rules, mods); err != nil {
		return "", CompileMeta{}, err
	}

	sortedMods := resolver.SortModules(mods)

	out, unresolved := render.Render(sortedMods, vars)

	for _, u := range unresolved {
		fmt.Fprintf(os.Stderr, "warning: unresolved variable: %s\n", u)
	}

	h := sha256.Sum256([]byte(out))
	hash := hex.EncodeToString(h[:])

	meta := CompileMeta{
		SelectedIDs:    selectedIDs,
		ClosureIDs:     closureIDs,
		Order:          order,
		Hash:           hash,
		UnresolvedVars: unresolved,
	}

	return out, meta, nil
}

func loadVarsFile(path string) (substitute.Vars, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var vars substitute.Vars
	if err := yaml.Unmarshal(data, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}

func buildSelectedIDs(opts CompileOptions) []string {
	selectedIDs := []string{
		"base",
		"modes/" + opts.Mode,
		"contracts/" + opts.Contract,
	}
	selectedIDs = append(selectedIDs, opts.Traits...)
	for _, g := range opts.Guardrails {
		selectedIDs = append(selectedIDs, "guardrails/"+g)
	}
	return selectedIDs
}

func buildModuleList(
	closureIDs []string,
	fromReq map[string]bool,
	selectedIDs []string,
	modByID map[string]*model.Module,
) ([]*model.Module, []string) {
	var mods []*model.Module
	var order []string

	for _, id := range closureIDs {
		m := modByID[id]
		m.FromReq = fromReq[id]
		m.Selected = resolver.Contains(selectedIDs, id)
		mods = append(mods, m)
		order = append(order, id)
	}

	return mods, order
}

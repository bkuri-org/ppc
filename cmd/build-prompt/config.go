package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	compilepkg "github.com/bkuri/ppc/internal/compile"
	"github.com/bkuri/ppc/internal/loader"
	"github.com/bkuri/ppc/internal/model"
	profilepkg "github.com/bkuri/ppc/internal/profile"
	"github.com/bkuri/ppc/internal/resolver"
)

type ResolvedConfig struct {
	Mode       string
	Contract   string
	Revisions  int
	Traits     []string
	Guardrails []string
	Vars       map[string]any
	VarsFile   string
	PromptsDir string
}

func NewResolvedConfigFromProfile(profileName string) (*ResolvedConfig, error) {
	profile, err := profilepkg.LoadProfile(profileName)
	if err != nil {
		return nil, err
	}

	cfg := ResolvedConfig{
		Mode:       profile.Mode,
		Contract:   profile.Contract,
		Revisions:  -1,
		Traits:     profile.Traits,
		Guardrails: []string{},
		Vars:       map[string]any{},
		VarsFile:   "",
		PromptsDir: "",
	}

	if profile.Revisions != nil {
		cfg.Revisions = *profile.Revisions
	}

	if profile.Vars != nil {
		cfg.Vars = profile.Vars
	}

	return &cfg, nil
}

func NewResolvedConfigFromDefaults(mode string, contract string) ResolvedConfig {
	return ResolvedConfig{
		Mode:       mode,
		Contract:   contract,
		Revisions:  -1,
		Traits:     []string{},
		Guardrails: []string{},
		Vars:       map[string]any{},
		VarsFile:   "",
		PromptsDir: "",
	}
}

func (c *ResolvedConfig) ApplyCLIOverrides(conservative, creative, terse, verbose *bool, revisions *int, contract, varsFile, guardrails *string) (*ResolvedConfig, error) {
	cfg := *c

	if conservative != nil && *conservative {
		cfg.Traits = append(cfg.Traits, "traits/conservative")
	}
	if creative != nil && *creative {
		cfg.Traits = append(cfg.Traits, "traits/creative")
	}
	if terse != nil && *terse {
		cfg.Traits = append(cfg.Traits, "traits/terse")
	}
	if verbose != nil && *verbose {
		cfg.Traits = append(cfg.Traits, "traits/verbose")
	}
	if revisions != nil {
		cfg.Revisions = *revisions
	}
	if contract != nil && *contract != "" {
		cfg.Contract = *contract
	}
	if varsFile != nil && *varsFile != "" {
		cfg.VarsFile = *varsFile
	}
	if guardrails != nil && *guardrails != "" {
		cfg.Guardrails = parseGuardrails(*guardrails, cfg.PromptsDir)
	}

	return &cfg, validateExclusiveGroups(cfg.Traits)
}

// parseGuardrails parses a comma-separated guardrail flag value.
// If value is "all", auto-discovers all guardrail modules in the prompts directory.
func parseGuardrails(value string, promptsDir string) []string {
	if value == "all" {
		return discoverAllGuardrails(promptsDir)
	}
	return parseCSV(value)
}

func discoverAllGuardrails(promptsDir string) []string {
	entries, err := os.ReadDir(filepath.Join(promptsDir, "guardrails"))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}
		out = append(out, strings.TrimSuffix(name, ".md"))
	}
	sort.Strings(out)
	return out
}

func (c *ResolvedConfig) ToCompileOptions() compilepkg.CompileOptions {
	return compilepkg.CompileOptions{
		Mode:       c.Mode,
		Contract:   c.Contract,
		Traits:     c.Traits,
		Guardrails: c.Guardrails,
		PromptsDir: c.PromptsDir,
		VarsFile:   c.VarsFile,
		Vars:       c.Vars,
	}
}

func validateExclusiveGroups(traits []string) error {
	modByID, errIf := loader.LoadModules("prompts")
	if errIf != nil {
		return fmt.Errorf("loading modules for exclusive group validation: %w", errIf)
	}

	rules, errIf := loader.LoadRules("prompts")
	if errIf != nil {
		return fmt.Errorf("loading rules for exclusive group validation: %w", errIf)
	}

	var mods []*model.Module
	for _, t := range traits {
		m, ok := modByID[t]
		if !ok {
			continue
		}
		mods = append(mods, m)
	}

	return resolver.ValidateExclusiveGroups(rules, mods)
}

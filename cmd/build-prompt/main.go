package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bkuri/ppc/internal/compile"
	"github.com/bkuri/ppc/internal/doctor"
	"github.com/bkuri/ppc/internal/lint"
	"github.com/bkuri/ppc/internal/loader"
)

// dief prints error to stderr and exits
func dief(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// explainOutput prints compilation metadata to stderr (CLI concern)
func explainOutput(meta compile.CompileMeta) {
	fmt.Fprintln(os.Stderr, "PPC explain")

	fmt.Fprintln(os.Stderr, "Selected IDs:")
	sel := append([]string{}, meta.SelectedIDs...)
	sort.Strings(sel)
	for _, id := range sel {
		fmt.Fprintf(os.Stderr, "  - %s\n", id)
	}

	fmt.Fprintln(os.Stderr, "Closure IDs (after requires):")
	cls := append([]string{}, meta.ClosureIDs...)
	sort.Strings(cls)
	for _, id := range cls {
		fmt.Fprintln(os.Stderr, "  - "+id)
	}

	fmt.Fprintln(os.Stderr, "Final order:")
	for _, id := range meta.Order {
		fmt.Fprintf(os.Stderr, "  - %s\n", id)
	}
}

func runExplore(args []string, promptsDir string) {
	fs := flag.NewFlagSet("explore", flag.ExitOnError)

	profile := fs.String("profile", "", "load preset configuration (e.g., ship)")
	conservative := fs.Bool("conservative", false, "include traits/conservative")
	creative := fs.Bool("creative", false, "include traits/creative")
	terse := fs.Bool("terse", false, "include traits/terse")
	verbose := fs.Bool("verbose", false, "include traits/verbose")
	revisions := fs.Int("revisions", -1, "revision budget (enables policies/revisions)")
	contract := fs.String("contract", "markdown", "contract module (code|markdown)")
	varsFile := fs.String("vars", "", "path to YAML file with variable definitions")
	outPath := fs.String("out", "", "write output to file")
	explain := fs.Bool("explain", false, "explain resolution steps to stderr")
	withHash := fs.Bool("hash", false, "prepend prompt-id hash header")
	proDir := fs.String("prompts", promptsDir, "prompts directory")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage:
  ppc explore [flags]

Explore mode generates a prompt for exploration tasks.

flags:`)
		fs.PrintDefaults()
	}

	fs.Parse(args)

	cfg := &ResolvedConfig{}
	if *profile != "" {
		profCfg, err := NewResolvedConfigFromProfile(*profile)
		if err != nil {
			dief("profile error: %v", err)
		}
		cfg = profCfg
	} else {
		defaults := NewResolvedConfigFromDefaults("explore", *contract)
		cfg = &defaults
	}

	cfg, err := cfg.ApplyCLIOverrides(conservative, creative, terse, verbose, revisions, contract, varsFile)
	if err != nil {
		dief("merge error: %v", err)
	}

	cfg.PromptsDir = *proDir

	opts := cfg.ToCompileOptions()

	out, meta, _ := compile.Compile(opts)

	if *withHash {
		out = fmt.Sprintf("<!-- prompt-id: sha256:%s -->\n\n%s", meta.Hash, out)
	}

	if *explain {
		explainOutput(meta)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(out), 0o644); err != nil {
			dief("failed to write %s: %v", *outPath, err)
		}
	}
	fmt.Print(out)
}

func runBuild(args []string, promptsDir string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)

	profile := fs.String("profile", "", "load preset configuration (e.g., ship)")
	conservative := fs.Bool("conservative", false, "include traits/conservative")
	creative := fs.Bool("creative", false, "include traits/creative")
	terse := fs.Bool("terse", false, "include traits/terse")
	verbose := fs.Bool("verbose", false, "include traits/verbose")
	revisions := fs.Int("revisions", -1, "revision budget (enables policies/revisions)")
	contract := fs.String("contract", "markdown", "contract module (code|markdown)")
	varsFile := fs.String("vars", "", "path to YAML file with variable definitions")
	outPath := fs.String("out", "", "write output to file")
	explain := fs.Bool("explain", false, "explain resolution steps to stderr")
	withHash := fs.Bool("hash", false, "prepend prompt-id hash header")
	proDir := fs.String("prompts", promptsDir, "prompts directory")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage:
  ppc build [flags]

Build mode generates a prompt for building/implementing features.

flags:`)
		fs.PrintDefaults()
	}

	fs.Parse(args)

	cfg := &ResolvedConfig{}
	if *profile != "" {
		profCfg, err := NewResolvedConfigFromProfile(*profile)
		if err != nil {
			dief("profile error: %v", err)
		}
		cfg = profCfg
	} else {
		defaults := NewResolvedConfigFromDefaults("build", *contract)
		cfg = &defaults
	}

	cfg, err := cfg.ApplyCLIOverrides(conservative, creative, terse, verbose, revisions, contract, varsFile)
	if err != nil {
		dief("merge error: %v", err)
	}

	cfg.PromptsDir = *proDir

	opts := cfg.ToCompileOptions()

	out, meta, _ := compile.Compile(opts)

	if *withHash {
		out = fmt.Sprintf("<!-- prompt-id: sha256:%s -->\n\n%s", meta.Hash, out)
	}

	if *explain {
		explainOutput(meta)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(out), 0o644); err != nil {
			dief("failed to write %s: %v", *outPath, err)
		}
	}
	fmt.Print(out)
}

func runShip(args []string, promptsDir string) {
	fs := flag.NewFlagSet("ship", flag.ExitOnError)

	profile := fs.String("profile", "", "load preset configuration (e.g., ship)")
	conservative := fs.Bool("conservative", false, "include traits/conservative")
	creative := fs.Bool("creative", false, "include traits/creative")
	terse := fs.Bool("terse", false, "include traits/terse")
	verbose := fs.Bool("verbose", false, "include traits/verbose")
	revisions := fs.Int("revisions", -1, "revision budget (enables policies/revisions)")
	contract := fs.String("contract", "markdown", "contract module (code|markdown)")
	varsFile := fs.String("vars", "", "path to YAML file with variable definitions")
	outPath := fs.String("out", "", "write output to file")
	explain := fs.Bool("explain", false, "explain resolution steps to stderr")
	withHash := fs.Bool("hash", false, "prepend prompt-id hash header")
	proDir := fs.String("prompts", promptsDir, "prompts directory")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage:
  ppc ship [flags]

Ship mode generates a prompt for release/deployment tasks.

flags:`)
		fs.PrintDefaults()
	}

	fs.Parse(args)

	cfg := &ResolvedConfig{}
	if *profile != "" {
		profCfg, err := NewResolvedConfigFromProfile(*profile)
		if err != nil {
			dief("profile error: %v", err)
		}
		cfg = profCfg
	} else {
		defaults := NewResolvedConfigFromDefaults("ship", *contract)
		cfg = &defaults
	}

	cfg, err := cfg.ApplyCLIOverrides(conservative, creative, terse, verbose, revisions, contract, varsFile)
	if err != nil {
		dief("merge error: %v", err)
	}

	cfg.PromptsDir = *proDir

	opts := cfg.ToCompileOptions()

	out, meta, _ := compile.Compile(opts)

	if *withHash {
		out = fmt.Sprintf("<!-- prompt-id: sha256:%s -->\n\n%s", meta.Hash, out)
	}

	if *explain {
		explainOutput(meta)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(out), 0o644); err != nil {
			dief("failed to write %s: %v", *outPath, err)
		}
	}
	fmt.Print(out)
}

func printGlobalUsage() {
	fmt.Fprintln(os.Stderr, `usage:
  ppc <subcommand> [flags]

 subcommands:
  explore    Generate prompt for exploration mode
  build      Generate prompt for build mode
  ship       Generate prompt for shipping mode
  doctor     Validate module structure and dependencies
  lint       Check prompt policies against lint rules

 global flags:
  --list     List all available modules
  --version  Show version information
  --help     Show this help message

 examples:
  ppc explore --conservative --revisions 1 --contract markdown
  ppc build --conservative --revisions 1 --contract code --explain
  ppc ship --creative --out AGENTS.md --hash
  ppc doctor --strict --json
  ppc lint --max-words 2000 --require-tags domain:*

 run 'ppc <subcommand> --help' for subcommand-specific options`)
}

func handleListModules(promptsDir string) {
	modByID, err := loader.LoadModules(promptsDir)
	if err != nil {
		dief("%v", err)
	}

	var ids []string
	for id := range modByID {
		ids = append(ids, id)
	}

	sort.Strings(ids)
	for _, id := range ids {
		m := modByID[id]
		desc := m.Front.Desc
		if desc == "" {
			desc = "(no desc)"
		}
		fmt.Printf("%-22s  %s\n", id, desc)
	}
}

func main() {
	// Handle global flags that don't require subcommand parsing
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		printVersion()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		printGlobalUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	// Default prompts directory
	promptsDir := "prompts"

	// Handle global meta-flags first
	if subcommand == "--list" {
		handleListModules(promptsDir)
		os.Exit(0)
	}

	if subcommand == "--help" || subcommand == "-h" || subcommand == "help" {
		printGlobalUsage()
		os.Exit(0)
	}

	// Dispatch to subcommand
	switch subcommand {
	case "explore":
		runExplore(args, promptsDir)
	case "build":
		runBuild(args, promptsDir)
	case "ship":
		runShip(args, promptsDir)
	case "doctor":
		fs := flag.NewFlagSet("doctor", flag.ExitOnError)
		strict := fs.Bool("strict", false, "treat warnings as errors")
		jsonOut := fs.Bool("json", false, "output machine-readable JSON")
		withStats := fs.Bool("stats", false, "include module statistics in JSON output")
		graphOut := fs.Bool("graph", false, "output Graphviz DOT format")
		outPath := fs.String("out", "", "write output to file")
		proDir := fs.String("prompts", promptsDir, "prompts directory")
		fs.Usage = func() {
			fmt.Fprintln(os.Stderr, `usage:
  ppc doctor [flags]

Checks module integrity, requires, cycles, and tag/rules sanity.

flags:`)
			fs.PrintDefaults()
		}
		fs.Parse(args)
		os.Exit(doctor.RunDoctor(*proDir, *strict, *jsonOut, *withStats, *graphOut, *outPath))

	case "lint":
		fs := flag.NewFlagSet("lint", flag.ExitOnError)
		maxWords := fs.Int("max-words", 0, "maximum total word count (0=disabled)")
		maxLines := fs.Int("max-lines", 0, "maximum total line count (0=disabled)")
		maxModules := fs.Int("max-modules", 0, "maximum number of modules (0=disabled)")
		maxModuleWords := fs.Int("max-module-words", 0, "maximum words per module (0=disabled)")
		maxDepth := fs.Int("max-depth", 0, "maximum dependency depth (0=disabled)")
		requireTags := fs.String("require-tags", "", "comma-separated list of required tags")
		forbidTags := fs.String("forbid-tags", "", "comma-separated list of forbidden tags")
		requireFields := fs.String("require-fields", "", "comma-separated list of required frontmatter fields")
		forbidEmptyBody := fs.Bool("forbid-empty-body", false, "fail if any module has empty body")
		jsonOut := fs.Bool("json", false, "output machine-readable JSON")
		proDir := fs.String("prompts", promptsDir, "prompts directory")
		fs.Usage = func() {
			fmt.Fprintln(os.Stderr, `usage:
  ppc lint [flags]

Checks prompt policies against configurable lint rules.

flags:`)
			fs.PrintDefaults()
		}
		fs.Parse(args)

		cfg := lint.Config{
			MaxWords:        *maxWords,
			MaxLines:        *maxLines,
			MaxModules:      *maxModules,
			MaxModuleWords:  *maxModuleWords,
			MaxDepth:        *maxDepth,
			RequireTags:     parseCSV(*requireTags),
			ForbidTags:      parseCSV(*forbidTags),
			RequireFields:   parseCSV(*requireFields),
			ForbidEmptyBody: *forbidEmptyBody,
		}

		result, err := lint.Run(*proDir, cfg)
		if err != nil {
			dief("lint error: %v", err)
		}

		if *jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(result); err != nil {
				dief("JSON encode error: %v", err)
			}
			if len(result.Violations) > 0 {
				os.Exit(2)
			}
			os.Exit(0)
		}

		if len(result.Violations) == 0 {
			fmt.Println("lint: OK")
			os.Exit(0)
		}

		fmt.Printf("lint: %d issue(s)\n", len(result.Violations))
		for _, v := range result.Violations {
			if v.Module != "" {
				fmt.Printf("  - [%s] %s: %s (%s)\n", v.Level, v.Rule, v.Message, v.Module)
			} else {
				fmt.Printf("  - [%s] %s: %s\n", v.Level, v.Rule, v.Message)
			}
		}
		os.Exit(2)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", subcommand)
		printGlobalUsage()
		os.Exit(1)
	}
}

// Package compile provides the core compilation API
package compile

type CompileOptions struct {
	Mode       string
	Contract   string
	Traits     []string
	Guardrails []string
	PromptsDir string
	VarsFile   string
	Vars       map[string]any
}

// CompileMeta provides metadata about the compilation
type CompileMeta struct {
	SelectedIDs    []string
	ClosureIDs     []string
	Order          []string
	Hash           string
	UnresolvedVars []string
}

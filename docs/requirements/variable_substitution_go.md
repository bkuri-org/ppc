# Implementation Plan: Variable Substitution (Go)

**Status:** Draft  
**Created:** 2026-03-11  
**Scope:** Minimal, deterministic Jinja2-style variable substitution for PPC

---

Nested variable substitution for `{{object.property}}` syntax. Variables load from external YAML/JSON files and substitute deterministically during compilation.

---

## Current State

PPC already has basic flat variable substitution:

```go
// internal/render/render.go
for k, v := range vars {
    body = strings.ReplaceAll(body, "{{"+k+"}}", v)
}
```

This works for `{{mode}}` but not `{{goals.target_return}}`.

---

## Design Goals

1. **Deterministic** - Same inputs → same output
2. **Explicit failure** - Warn on unresolved variables, don't silently drop
3. **No templating logic** - Only substitution, no `{% if %}` blocks
4. **External config** - Variables loaded from files, not embedded

---

## Architecture

### 1. Variable Sources

Variables come from a single YAML/JSON file specified at compile time:

```
ppc explore --vars ~/.config/ppc/vars.yml
```

**File format (YAML):**

```yaml
goals:
  primary_objective: growth
  target_annual_return: 15.0
  time_horizon_months: 12
  priority_metric: sharpe

user:
  risk_style: balanced

safety:
  max_drawdown: 0.20
```

**Rationale:** Single file is simpler than multiple sources. Users can merge configs externally if needed.

### 2. Substitution Engine

**New package:** `internal/substitute/substitute.go`

```go
package substitute

// Vars represents nested variable structure
type Vars map[string]any

// Substitute replaces {{path.to.value}} with actual values
// Supports nested access via dot notation
func Substitute(content string, vars Vars) string

// ResolvePath resolves "goals.target_return" to the actual value
// Returns (value, found)
func ResolvePath(vars Vars, path string) (any, bool)
```

**Implementation approach:**

1. Find all `{{...}}` patterns using regex
2. Extract path (e.g., `goals.target_return`)
3. Resolve path against vars map
4. Replace placeholder with stringified value
5. Log warning for unresolved variables (don't fail)

### 3. Modified Render

Update `internal/render/render.go`:

```go
func Render(mods []*model.Module, vars substitute.Vars) string {
    var b strings.Builder
    for i, m := range mods {
        if i > 0 {
            b.WriteString("\n\n")
        }
        body := substitute.Substitute(m.Body, vars)
        b.WriteString(strings.TrimRight(body, "\n"))
    }
    return strings.TrimRight(b.String(), "\n") + "\n"
}
```

### 4. CLI Integration

Add `--vars` flag to all subcommands:

```go
varsFile := fs.String("vars", "", "YAML/JSON file with variable definitions")
```

**Loading logic:**

```go
func loadVars(path string) (substitute.Vars, error) {
    if path == "" {
        return substitute.Vars{}, nil
    }
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
```

### 5. Cache Key Update

The hash already includes module content. Add vars file hash:

```go
func computeHash(out string, varsPath string) string {
    h := sha256.New()
    h.Write([]byte(out))
    if varsPath != "" {
        varsData, _ := os.ReadFile(varsPath)
        h.Write(varsData)
    }
    return hex.EncodeToString(h.Sum(nil))
}
```

---

## File Changes

| File | Change |
|------|--------|
| `internal/substitute/substitute.go` | **NEW** - Substitution engine |
| `internal/substitute/substitute_test.go` | **NEW** - Unit tests |
| `internal/render/render.go` | Replace flat substitution with nested |
| `internal/compile/types.go` | Add `VarsFile string` to `CompileOptions` |
| `internal/compile/compile.go` | Load vars file, pass to render |
| `cmd/build-prompt/main.go` | Add `--vars` flag to all subcommands |

---

## Implementation Steps

### Phase 1: Core Substitution

1. Create `internal/substitute/substitute.go`
2. Implement `ResolvePath()` with dot-notation support
3. Implement `Substitute()` with regex matching
4. Handle type conversion (int/float → string)
5. Log warnings for unresolved variables

### Phase 2: Integration

1. Update `CompileOptions` with `VarsFile`
2. Add vars loading in `compile.Compile()`
3. Update `Render()` signature and call
4. Add `--vars` flag to CLI

### Phase 3: Testing

1. Unit tests for substitution engine
2. Integration test with sample vars file
3. Golden test updates for new functionality

---

## Edge Cases

| Case | Behavior |
|------|----------|
| Missing vars file | Error, exit 2 |
| Empty vars file | Valid, no substitution |
| `{{unknown.path}}` | Log warning, keep placeholder |
| `{{goals}}` (object) | JSON stringify |
| `{{goals.null_value}}` | Replace with empty string |
| Circular references | N/A (no templating logic) |

---

## Example Usage

**vars.yml:**

```yaml
goals:
  target_return: 20.0
  horizon: 24
```

**prompts/policies/outcome_targets.md:**

```markdown
---
id: policies/outcome_targets
desc: User-defined outcome targets
priority: 2
tags: [policy:outcomes]
requires: [base]
---

## Outcome Targets

Target return: {{goals.target_return}}%
Time horizon: {{goals.horizon}} months
```

**Command:**

```bash
ppc explore --vars vars.yml --explain
```

**Output includes:**

```markdown
## Outcome Targets

Target return: 20.0%
Time horizon: 24 months
```

---

## Out of Scope (Phase 2+)

- `{% if %}` conditional blocks
- Filters like `{{value | round(2)}}`
- Multiple vars files
- Default values syntax (`{{var | default}}`)
- Profile-level vars file references

---

## Acceptance Criteria

1. `--vars` loads YAML/JSON config
2. `{{object.property}}` resolves correctly
3. Unresolved variables log warnings but don't fail
4. Empty vars file produces unchanged output
5. Cache hash includes vars file content
6. All existing tests pass
7. New tests cover substitution logic

---

## Risks

| Risk | Mitigation |
|------|------------|
| Complex nested structures | Limit depth, document limits |
| Regex performance | Use compiled regex, benchmark |
| Type formatting ambiguity | Always stringify, let user format in prompt |

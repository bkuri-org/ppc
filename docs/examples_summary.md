# PPC Examples — Implementation Summary

All 5 examples from `docs/examples_prd.md` have been successfully created and validated.

## Example Overview

| # | Name | Complexity | Modules | Status |
|---|------|------------|--------|
| 01 | Basic Prompt Composition | 1 | 5 | Validated |
| 02 | Team Style Guide Policy | 2 | 7 | Validated |
| 03 | Knowledge Sharing Policy | 3 | 7 | Validated |
| 04 | Product PRD Review Flow | 4 | 7 | Validated |
| 05 | RAG Governance Policy | 5 | 7 | Validated |

## Features Demonstrated

### Example 01 — Basic Prompt Composition
- PPC as prompt compiler (replaces giant prompt files)
- Modular behavior composition (base, mode, traits, contracts)
- Deterministic ordering
- Exclusive tag groups (risk, output)
- Profile-based execution

**Profiles:**
- `explore-creative` — exploration with creative trait
- `explore-conservative` — exploration with conservative trait

### Example 02 — Team Style Guide Policy
- Organization-wide behavioral policy
- Tone enforcement (terse vs verbose)
- Formatting standards
- Exclusive tag groups (tone, output)

**Profiles:**
- `default` — terse + markdown
- `design-reviews` — verbose + markdown

### Example 03 — Knowledge Sharing Policy
- Conversational governance
- 7-day conversation windows
- Requirements traceability (PRD citations)
- RAG source attribution
- Contracts as process definition

**Profiles:**
- `knowledge-sharing` — complete governance setup

### Example 04 — Product PRD Review Flow (Variables Introduced Here)
- Multi-stage review (exploration, definition, ship)
- Artifact-driven development
- Variable substitution (`{{product_name}}`, `{{target_user}}`, `{{risk_level}}`)
- Heavy contract module (60-80 lines)
- Acceptance criteria, non-goals, risk assessment

**Profiles:**
- `explore` — NewFeature, Developers, low risk
- `ship` — ProductionRelease, Customers, high risk

### Example 05 — RAG Governance Policy
- Enterprise AI governance
- Multiple exclusive groups (audience, output)
- Source ranking (5-level hierarchy)
- Citation format requirements (`[source-id:confidence:timestamp]`)
- Escalation rules (mandatory vs conditional)
- Heavy contract module (60-80 lines)

**Profiles:**
- `internal` — full access, moderate escalation
- `client-facing` — limited access, strict escalation

## Directory Structure

Each example is self-contained with:

```
example-name/
  README.md                    # Problem → Solution → Value template
  prompts/
    base.md                   # Example-specific identity
    modes/
      explore.md                # Exploration mode
    traits/                     # 02 only
      terse.md
      verbose.md
    policies/                   # 03, 04, 05
      [various policies]
    contracts/
      markdown.md               # All examples
      [specific contracts]       # 03, 04, 05
    rules.yml                  # Exclusive group definitions
  profiles/
    [profile configs]
```

## README Specifications

All READMEs follow the template:

1. **Problem**: Real-world pain point
2. **Why this exists**: Why ad-hoc approaches fail
3. **What this example demonstrates**: Concepts introduced
4. **How to run**: `ppc doctor` and profile commands
5. **Output**: Sample compiled output (truncated, 30-60 lines)
6. **What to copy into your project**: Copy-paste guidance
7. **Common failure modes**: Tag conflicts, missing variables, citation errors
8. **CI**: Reference to canonical workflow
9. **Key takeaway**: One-sentence insight

## Validation Results

All examples pass `ppc doctor`:

```
=== 01-basic-prompt ===
doctor: OK (5 modules)

=== 02-team-style-guide ===
doctor: OK (7 modules)

=== 03-knowledge-sharing-policy ===
doctor: OK (7 modules)

=== 04-product-prd-review ===
doctor: OK (7 modules)

=== 05-rag-governance-policy ===
doctor: OK (7 modules)
```

Note: "unreachable modules" warnings are expected and intentional — traits/policies are selected via profiles, not via requires.

## Testing Commands

Validate each example:

```bash
# Example 01
cd examples/01-basic-prompt
ppc doctor
ppc explore --profile explore-creative
ppc explore --profile explore-conservative

# Example 02
cd examples/02-team-style-guide
ppc doctor
ppc explore --profile default
ppc explore --profile design-reviews

# Example 03
cd examples/03-knowledge-sharing-policy
ppc doctor
ppc explore --profile knowledge-sharing

# Example 04 (variables)
cd examples/04-product-prd-review
ppc doctor
ppc explore --profile explore  # Substitutes: NewFeature, Developers, low
ppc explore --profile ship    # Substitutes: ProductionRelease, Customers, high

# Example 05 (multiple exclusive groups)
cd examples/05-rag-governance-policy
ppc doctor
ppc explore --profile internal         # Audience: internal
ppc explore --profile client-facing    # Audience: client-facing
```

## Module Content Style

All modules follow realistic, AGENTS.md-level content:

- **20-30 lines** for standard modules
- **60-80 lines** for heavy modules (Example 04 contracts/prd-review.md, Example 05 contracts/rag-governance.md)
- **No simplification** — copy-pasteable for production
- **Explicit constraints** — "You must not" and "You do" sections
- **Anti-patterns** — Clear failure cases

## Success Criteria

- Users recognize their own problems immediately
- Examples are fully self-contained
- Each example demonstrates exactly one new conceptual idea
- Complexity ramps progressively (1 to 5)
- READMEs follow problem → solution → value structure
- All examples pass `ppc doctor`
- Variable substitution demonstrated (Example 04)
- Multiple exclusive groups demonstrated (Example 05)
- Sample output included (truncated)
- Common failure modes documented
- CI integration path clear

## Next Steps

The examples are ready for:

1. **Review**: Validate content matches organizational needs
2. **Testing**: Run each profile and verify compiled output
3. **Documentation**: Update main README to reference examples
4. **CI Integration**: Ensure `examples/workflows/prompt-lint.yml` can validate all examples
5. **User Testing**: Get feedback on clarity and copy-paste utility

## Implementation Notes

- **Total files created**: 55+ files across 5 examples
- **Total lines of content**: 2000+ lines
- **Average README length**: 137-206 lines
- **Heavy modules**: 2 (prd-review, rag-governance)
- **Variable examples**: 1 (product-prd-review)
- **Multiple exclusive groups**: 1 (rag-governance-policy)

## Compatibility

All examples are compatible with:
- PPC v0.2.0 (current version)
- Standard Markdown parsers
- CommonMark syntax
- YAML frontmatter
- Variable substitution (`{{varname}}`)

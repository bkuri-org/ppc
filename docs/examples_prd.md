# PPC Examples — Product Requirements Document (PRD)

## Purpose

The `examples/` directory exists to **teach PPC by demonstration**.

Rather than explaining concepts abstractly, examples provide:

- concrete, runnable prompt-policy repositories
- realistic use cases users already recognize
- increasing complexity through composition
- copy-paste starting points for real projects

Examples are reference implementations, not documentation.

---

## Design Goals

Each example must:

- be complete and runnable
- compile successfully with `ppc doctor --strict`
- include a short README explaining *why* it exists
- demonstrate exactly one new conceptual idea
- build on concepts introduced earlier

Non‑goals:

- no artificial demos
- no toy prompts
- no placeholder text
- no fake complexity

Each example should feel like something a real team would actually use.

---

## Directory Structure

```
examples/
  01-basic-prompt/
  02-team-style-guide/
  03-knowledge-sharing-policy/
  04-product-prd-review/
  05-rag-governance-policy/
```

Each example is self‑contained.

```
example-name/
  README.md
  prompts/
  profiles/
  rules.yml
```

Users should be able to run:

```bash
ppc doctor
ppc explore --profile <profile>
```

from inside any example directory.

---

# Example 01 — Basic Prompt Composition

### Folder

```
01-basic-prompt/
```

### Concept Introduced

- PPC as a prompt compiler
- modules
- deterministic ordering

### Demonstrates

- `base`
- `modes/explore`
- simple traits
- markdown output contract

### README should explain

- why prompts drift
- how modules compose
- how output is deterministic

### Expected modules

- base
- modes/explore
- traits/creative
- traits/conservative
- contracts/markdown

### Key takeaway

> "This is what PPC replaces: one giant prompt file."

---

# Example 02 — Team Style Guide Policy

### Folder

```
02-team-style-guide/
```

### Concept Introduced

- organization‑wide behavioral policy

### Demonstrates

- tone enforcement
- formatting standards
- output rules
- exclusive tag groups

### New ideas

- `tone:*` exclusive group
- company voice enforcement

### Example modules

- traits/terse
- traits/verbose
- policies/style-guide
- policies/formatting

### Profiles

```yaml
profiles/default.yml
profiles/verbose.yml
```

### Key takeaway

> "PPC can enforce team-wide communication standards."

---

# Example 03 — Knowledge Sharing Policy

### Folder

```
03-knowledge-sharing-policy/
```

### Concept Introduced

- conversational governance
- requirements traceability

### Demonstrates

- contracts as process definition
- conversation period enforcement
- PRD/RAG evidence

### Core modules

- contracts/knowledge-sharing
- policies/conversation-period
- policies/requirements-evidence
- policies/rag-citations

### Profiles

```yaml
profiles/knowledge-sharing.yml
```

### README focus

Explain:

- why conversations need structure
- how PPC standardizes review discussions
- how CI can gate artifact presence

### Key takeaway

> "PPC can encode organizational process — not just prompt text."

---

# Example 04 — Product PRD Review Flow

### Folder

```
04-product-prd-review/
```

### Concept Introduced

- multi-stage review
- artifact-driven development

### Demonstrates

- PRD review contract
- acceptance criteria
- non-goals enforcement
- revision budgeting

### Example modules

- contracts/prd-review
- policies/acceptance-criteria
- policies/non-goals
- policies/risk-assessment

### Profiles

```yaml
profiles/explore.yml
profiles/ship.yml
```

### README focus

- treating PRDs as first-class artifacts
- preventing scope creep
- separating exploration from approval

### Key takeaway

> "PPC can formalize how products are designed, not just how prompts are written."

---

# Example 05 — RAG Governance Policy (Advanced)

### Folder

```
05-rag-governance-policy/
```

### Concept Introduced

- high-risk knowledge systems
- auditability
n### Demonstrates

- mandatory citation policy
- hallucination containment
- retrieval transparency
- escalation rules

### Example modules

- contracts/rag-governance
- policies/source-ranking
- policies/citation-format
- policies/failure-escalation

### Profiles

```yaml
profiles/internal.yml
profiles/client-facing.yml
```

### Advanced concepts

- multiple contracts
- shared policies
- exclusive audience groups

### README focus

- why RAG needs governance
- legal/compliance alignment
- audit-ready AI behavior

### Key takeaway

> "PPC can define AI governance policies suitable for regulated environments."

---

## Complexity Ramp

| Example | Complexity | New Concept |
|------|------|------|
| 01 | 1 | Composition |
| 02 | 2 | Policy enforcement |
| 03 | 3 | Process governance |
| 04 | 4 | Product workflows |
| 05 | 5 | Enterprise AI governance |

Each example assumes mastery of the previous one.

---

## Documentation Rules

Each example README must include:

1. **Problem statement**
2. **Why PPC helps**
3. **Folder layout**
4. **How to run it**
5. **What to copy into your own repo**

No README should exceed ~2 pages.

---

## Success Criteria

The examples folder succeeds if:

- users recognize their own problems immediately
- examples are frequently copied verbatim
- most onboarding questions are answered by pointing to one example
- PPC’s purpose becomes obvious without reading the code

---

## Measure of Success

Examples succeed when users recognize their own problems and can adapt an example for their use case. If an example doesn't map to real work, it doesn't belong here.

---


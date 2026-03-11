# PPC Modes

Modes guide agent behavior through the software development lifecycle.

## When to Use Each Mode

| Mode | Phase | Question It Answers |
|------|-------|---------------------|
| `explore` | Discovery | "What should we do?" |
| `build` | Implementation | "How do we do it?" |
| `ship` | Release | "Is it safe to ship?" |

---

## Explore

**Philosophy:** Breadth before depth. Generate options. Recommend.

### Use When
- Investigating an unfamiliar codebase
- Answering "how should I..." questions
- Evaluating tradeoffs between approaches
- Researching without making changes

### Avoid When
- You already know the solution
- The task is straightforward implementation

### Compiled From
- `prompts/base.md`
- `prompts/modes/explore.md`
- Selected traits (conservative, creative, terse, verbose)
- Selected contract (markdown, code)

### Example
```bash
ppc explore --creative --vars ./goals.yaml
```

---

## Build

**Philosophy:** Plan then execute. Minimal moving parts. Readable code.

### Use When
- Implementing a feature
- Fixing a bug
- Refactoring code
- Executing a known plan

### Avoid When
- You don't know what to build yet (use explore)
- You're preparing for production release (use ship)

### Compiled From
- `prompts/base.md`
- `prompts/modes/build.md`
- Selected traits and contracts

### Example
```bash
ppc build --conservative --contract code --revisions 2
```

---

## Ship

**Philosophy:** Stability first. Boring solutions. Explicit failure.

### Use When
- Preparing for release or deployment
- Making low-risk changes to production
- Finalizing work
- Code review for merge

### Avoid When
- You need to experiment (use explore)
- You're building new features (use build)

### Compiled From
- `prompts/base.md`
- `prompts/modes/ship.md`
- Selected traits and contracts
- Automatically includes `risk:low` tag

### Example
```bash
ppc ship --terse --hash --out AGENTS.md
```

---

## Mode Selection Guide

| Question | Mode |
|----------|------|
| "What's the best way to...?" | explore |
| "Implement this feature" | build |
| "Deploy to production" | ship |
| "Fix this bug" | build |
| "Should I use X or Y?" | explore |
| "Review this PR" | explore |
| "Is this ready to merge?" | ship |

---

## Combining Modes in Workflows

Real-world workflows often combine modes:

### Feature Development
```bash
# 1. Research the approach
ppc explore --creative --out research.md

# 2. Implement the feature
ppc build --conservative --contract code

# 3. Prepare for merge
ppc ship --terse --hash
```

### Bug Investigation
```bash
# 1. Understand the problem
ppc explore --conservative

# 2. Fix it
ppc build --conservative --contract code --revisions 1
```

### Production Hotfix
```bash
# Skip explore, go straight to fix
ppc build --conservative --contract code

# Ship with maximum caution
ppc ship --terse --hash
```

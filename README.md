# PPC — Prompt Policy Compiler (v0.3.1)

PPC compiles small Markdown behavior modules into a single deterministic prompt (stdout-first).

**New to PPC?** Start with [Examples](#examples) below to see real-world use cases.

## Quick Start

```bash
# Install
go install github.com/bkuri/ppc/cmd/build-prompt@v0.3.1

# Compile a prompt
ppc explore --creative --out my-prompt.md

# Validate your prompt repository
ppc doctor --strict

# See available modules
ppc --list
```

## Installation

### Method 1: Download from Releases (Recommended)

1. Download for your platform:
   ```bash
   curl -fsSL -o ppc.tar.gz \
     https://github.com/bkuri/ppc/releases/download/v0.3.1/ppc_v0.3.1_linux_amd64.tar.gz
   ```

2. Extract and install:
   ```bash
   tar -xzf ppc.tar.gz
   chmod +x linux_amd64/ppc
   sudo mv linux_amd64/ppc /usr/local/bin/ppc
   ```

3. Verify:
   ```bash
   ppc --version
   ```

Available platforms: `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64`, `windows_amd64`.

### Method 2: Verify Checksums

Always verify release integrity:

```bash
# Download checksums
curl -fsSL -O https://github.com/bkuri/ppc/releases/download/v0.3.1/checksums.txt

# Verify your downloaded archive
sha256sum -c --ignore-missing checksums.txt
```

Should output: `ppc_v0.3.1_linux_amd64.tar.gz: OK`

See [docs/verification.md](docs/verification.md) for detailed verification guide.

### Method 3: Go Install

Pin to specific version for reproducibility:

```bash
go install github.com/bkuri/ppc/cmd/build-prompt@v0.3.1
```

The binary installs to `$GOPATH/bin/ppc` (usually `~/go/bin/ppc`).

### Method 4: Build from Source

```bash
git clone https://github.com/bkuri/ppc.git
cd ppc
go build -o ppc ./cmd/build-prompt
```

## Versioning

PPC follows semantic versioning: `vX.Y.Z`

- X: Major version (breaking changes)
- Y: Minor version (new features, backward compatible)
- Z: Patch version (bug fixes)

To install a specific version, use:
```bash
go install github.com/bkuri/ppc/cmd/build-prompt@v0.3.1
```

Check for latest releases at:
https://github.com/bkuri/ppc/releases

## Build

```bash
go build -o ppc ./cmd/build-prompt
```

## Examples

If you want to understand PPC quickly, start here:

All examples are complete, runnable prompt-policy repositories that increase in complexity.

| # | Name | Complexity | What You'll Learn |
|---|------|-----------|-------------------|
| 01 | [Basic Prompt Composition](examples/01-basic-prompt) | 1 | Modular composition, deterministic ordering |
| 02 | [Team Style Guide](examples/02-team-style-guide) | 2 | Policy enforcement, exclusive groups |
| 03 | [Knowledge Sharing Policy](examples/03-knowledge-sharing-policy) | 3 | Process governance, traceability |
| 04 | [Product PRD Review](examples/04-product-prd-review) | 4 | Multi-stage workflows, variable substitution |
| 05 | [RAG Governance](examples/05-rag-governance-policy) | 5 | Enterprise governance, multiple exclusive groups |

Each example includes a README explaining the problem it solves and how to run it.

To try an example:

```bash
cd examples/01-basic-prompt
ppc doctor                    # Validate structure
ppc explore --profile explore-creative | head -20
```

## Modes

PPC uses three modes aligned with the software lifecycle:

| Mode | Use When | One-liner |
|------|----------|-----------|
| `explore` | Researching, recommending, answering questions | "What's the best approach?" |
| `build` | Implementing, coding, making changes | "Make it work." |
| `ship` | Releasing, deploying, finalizing | "Don't break anything." |

See [docs/modes.md](docs/modes.md) for detailed mode guide and selection flowchart.

## Subcommands

PPC uses a subcommand-based CLI:

```bash
ppc <subcommand> [flags]
```

### Mode Subcommands

Generate prompts for specific modes:

```bash
./ppc explore --conservative --revisions 1 --contract markdown
./ppc build --conservative --revisions 1 --contract code --explain
./ppc ship --creative --out AGENTS.md --hash
```

### Doctor Subcommand

Validate module structure and dependencies:

```bash
./ppc doctor --strict
./ppc doctor --json
./ppc doctor --prompts custom/
```

### Global Flags

```bash
./ppc --list        # List all available modules
./ppc --help        # Show help
```

### Flags Per Subcommand

Each mode subcommand supports:

```
--conservative          Include conservative trait (boring, stable)
--creative              Include creative trait (novelty encouraged)
--terse                 Include terse trait (brief, concise)
--verbose               Include verbose trait (detailed, expansive)
--revisions N           Enable policies/revisions with budget N
--contract TYPE         Output contract: code|markdown (default: markdown)
--out PATH              Write output to file (default: stdout)
--explain               Print resolution details to stderr
--hash                  Prepend SHA256 prompt-id header
--prompts DIR           Prompts directory (default: prompts)
--vars PATH             YAML/JSON file with variable definitions
```

## Layout

- `prompts/` contains Markdown modules with optional YAML frontmatter.
- `prompts/rules.yml` defines `exclusive_groups` for keyed tags (`group:value`).

Deterministic: same inputs produce same output. Fails loudly on missing modules, tag conflicts, and circular requires.

## Documentation

- **[PRD.md](PRD.md)** — Product requirements and design principles
- **[ROADMAP.md](ROADMAP.md)** — Long-term vision and planned features
- **[CHANGELOG.md](CHANGELOG.md)** — Version history and changes
- **[CONTRIBUTING.md](CONTRIBUTING.md)** — How to contribute
- **[docs/github-actions.md](docs/github-actions.md)** — CI/CD setup
- **[docs/verification.md](docs/verification.md)** — Checksum verification

## Support & Community

- **Found a bug?** [Open an issue](https://github.com/bkuri/ppc/issues)
- **Have a question?** [Start a discussion](https://github.com/bkuri/ppc/discussions)
- **Want to contribute?** See [CONTRIBUTING.md](CONTRIBUTING.md)
- **Code of Conduct** — See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## CI Integration

This repository includes automated workflows:

- **[lint.yml](.github/workflows/lint.yml)** — Tests and code quality
- **[validate-examples.yml](.github/workflows/validate-examples.yml)** — Example validation
- **[release.yml](.github/workflows/release.yml)** — Automated releases

See [docs/github-actions.md](docs/github-actions.md) for integration into your own repository.

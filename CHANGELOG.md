# Changelog

All notable changes to PPC are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.2] — 2026-03-12

### Added

- **Tests**: Comprehensive test coverage for new features
  - `internal/lint`: 93.7% coverage (11 new tests)
  - `internal/doctor`: Full coverage for validation logic

### Fixed

- Fixed graph tests to use correct prompts directory path

## [0.3.1] — 2026-03-11

### Added

- **Manpage**: Full CLI documentation in `man/ppc.1`

## [0.3.0] — 2026-03-11

### Added

- **Lint Command**: New `ppc lint` subcommand for prompt policy validation
  - `--forbid-empty-body`: Fail if any module has empty body
  - `--forbid-tags`: Comma-separated list of forbidden tags
  - `--max-depth`: Maximum dependency depth
  - `--max-lines`: Maximum total line count
  - `--max-module-words`: Maximum words per module
  - `--max-modules`: Maximum number of modules
  - `--json`: Machine-readable JSON output

- **Nested Variable Substitution**: Support for `--vars` flag with nested variable references

- **Modes Documentation**: Detailed mode guide in `docs/modes.md`

## [0.2.0] — 2026-01-23

### Added

- **Examples**: Five complete, runnable reference implementations demonstrating PPC capabilities
  - Example 01: Basic Prompt Composition (modular composition, deterministic ordering)
  - Example 02: Team Style Guide Policy (policy enforcement, tone groups)
  - Example 03: Knowledge Sharing Policy (process governance, traceability)
  - Example 04: Product PRD Review Flow (multi-stage workflows, variable substitution)
  - Example 05: RAG Governance Policy (enterprise governance, multiple exclusive groups)

- **Documentation**: Comprehensive user and contributor guides
  - `docs/examples_prd.md`: Specification for example implementations
  - `docs/examples_summary.md`: Implementation summary of all 5 examples
  - `CONTRIBUTING.md`: Contribution guidelines for maintainers and community
  - `CHANGELOG.md`: Version history (this file)
  - Enhanced README with quick-start and example links

- **CI/CD**: GitHub Actions workflows for quality assurance
  - `lint.yml`: Automated testing and code validation
  - `validate-examples.yml`: Example compilation and validation
  - `release.yml`: Automated release artifact building (previously existed)

- **Community Infrastructure**
  - CODE_OF_CONDUCT.md: Community standards
  - .github/ISSUE_TEMPLATE/bug_report.md: Bug report template
  - .github/PULL_REQUEST_TEMPLATE.md: PR submission template

### Changed

- **Documentation**: Refined README.md with better navigation and example references
- **Cleanup**: Removed temporary phase documentation files from docs/

### Fixed

- None (all fixes are in earlier releases)

### Deprecated

- None (stable release)

## [0.1.0] — 2025-12-15

### Added

**Core Compiler**
- Markdown module loading and parsing
- YAML frontmatter extraction (id, desc, priority, tags, requires)
- Deterministic module ordering (by layer and priority)
- Transitive requires expansion
- Circular dependency detection
- Tag validation and exclusive group enforcement
- Variable substitution (`{{varname}}`)
- Profile-based configuration

**CLI**
- Three mode subcommands: `explore`, `build`, `ship`
- Module selection via flags (traits, contracts, policies)
- Output to stdout or file (`--out`)
- Prompt hashing (`--hash`)
- Explain mode (`--explain`) for debugging
- Module listing (`--list`)

**Validation**
- `ppc doctor` for repository linting
- Strict mode (`--strict`) treating warnings as errors
- JSON output (`--json`) for machine parsing
- Graph visualization (`--graph`) for dependency analysis
- Statistics reporting (`--stats`)

**Configuration**
- `rules.yml` for exclusive group definitions
- `profiles/` directory for preset configurations
- Support for multiple mode/contract/trait combinations

**Testing**
- Golden snapshot tests for output consistency
- Basic compilation tests
- Graph and ordering tests
- Integration tests with example repositories

**Documentation**
- README.md with installation, usage, and layout
- PRD.md with detailed product specification
- ROADMAP.md with long-term vision (Phases 0-4)
- docs/github-actions.md with CI setup guide
- docs/verification.md with checksum verification

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and improvements.

## Version Naming

- **v0.1.x**: Baseline feature completeness
- **v0.2.x**: Examples and community infrastructure
- **v1.0.0**: Stable, production-ready release
- **v1.x+**: Feature additions and refinements

## Security

For security vulnerabilities, please email the maintainer privately rather than opening a public issue.

See [docs/verification.md](docs/verification.md) for checksum verification procedures.

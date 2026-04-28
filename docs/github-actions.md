# GitHub Actions for PPC

Automatically validate prompt modules in CI/CD pipelines.

---

## Download from Releases (Recommended)

**Advantages:**
- No Go toolchain required in CI
- Fast, deterministic version pinning
- Works offline if binary is cached

**Workflow:**

\`\`\`yaml
name: Prompt Policy Check

on: [push, pull_request]

jobs:
  doctor:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install ppc
        run: |
          set -euo pipefail
          VERSION=v0.3.1
          URL="https://github.com/bkuri/ppc/releases/download/${VERSION}/ppc_${VERSION}_linux_amd64.tar.gz"
          curl -fsSL -o ppc.tar.gz "$URL"
          tar -xzf ppc.tar.gz
          chmod +x ppc
          ./ppc --version

      - name: Doctor
        run: ./ppc doctor --strict --json --out report.json

      - uses: actions/upload-artifact@v4
        with:
          name: doctor-report
          path: report.json
\`\`\`

---

## Build from Source (Alternative)

**Advantages:**
- No release binary required
- Always builds from HEAD
- Simpler setup

**Workflow:**

\`\`\`yaml
name: Prompt Policy Check

on: [push, pull_request]

jobs:
  doctor:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - uses: actions/checkout@v4

      - name: Build ppc
        run: go build -o ppc ./cmd/build-prompt

      - name: Doctor
        run: ./ppc doctor --strict --json --out report.json

      - uses: actions/upload-artifact@v4
        with:
          name: doctor-report
          path: report.json
\`\`\`

---

## Optional Graph Integration

Add dependency graph visualization to your artifacts:

\`\`\`yaml
      # Optional: add graph generation
      - name: Generate Graph
        run: ./ppc doctor --graph --out deps.dot

      - name: Upload Graph Artifact
        uses: actions/upload-artifact@v4
        with:
          name: doctor-graph
          path: deps.dot
\`\`\`

Then manually:
1. Download \`deps.dot\` from artifacts
2. Render with \`dot -Tpng deps.dot -o deps.png\`
3. View dependency structure

---

## How to Pin Versions

### Pin to release version

\`\`\`yaml
      - name: Install ppc
        env:
          VERSION: v0.3.1  # Pin to specific release
\`\`\`

Update when upgrading PPC.

### Pin to commit (build from source)

\`\`\`yaml
      - uses: actions/checkout@v4
        with:
          ref: abc1234  # Specific commit SHA
\`\`\`

### Auto-update (not recommended)

Replace \`v0.3.1\` with a branch name like \`main\`, but this sacrifices determinism.

---

## Example jq Checks

Extract and validate report status:

\`\`\`yaml
      - name: Validate Results
        run: |
          STATUS=$(jq -r .status report.json)
          if [ "$STATUS" != "ok" ]; then
            echo "Prompt policy check failed"
            cat report.json
            exit 1
          fi
\`\`\`

---

## Troubleshooting

**Binary not found:** Verify release asset name matches your platform (\`linux_amd64\`, \`darwin_arm64\`, etc.)

**Doctor fails:** Check that your \`prompts/\` directory exists and contains valid modules.

**Graph not rendering:** \`dot\` is not required for CI; download artifact and render locally.

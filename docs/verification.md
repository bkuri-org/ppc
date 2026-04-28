# Verifying PPC Release Integrity

## Why Verify?

Release verification ensures:
- Downloaded file is uncorrupted
- No tampering occurred during transfer
- File matches published checksum

## Quick Verification

1. Download release and checksums:
   ```bash
   VERSION=v0.3.1
   curl -fsSL -O https://github.com/bkuri/ppc/releases/download/${VERSION}/ppc_${VERSION}_linux_amd64.tar.gz
   curl -fsSL -O https://github.com/bkuri/ppc/releases/download/${VERSION}/checksums.txt
   ```

2. Verify checksum:
   ```bash
   sha256sum -c --ignore-missing checksums.txt
   ```

Expected output: `ppc_v0.3.1_linux_amd64.tar.gz: OK`

## Manual Verification

If `sha256sum` not available:

1. View checksums:
   ```bash
   cat checksums.txt
   ```

2. Compute and compare:
   ```bash
   sha256sum ppc_v0.3.1_linux_amd64.tar.gz
   ```

3. Compare output with checksums.txt

## Verification Failures

If verification fails:
1. Delete corrupted download
2. Re-download archive
3. Verify checksums again
4. Report issue: https://github.com/bkuri/ppc/issues



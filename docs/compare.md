# envault compare

Compare two encrypted vault files and show key differences.

## Usage

```
envault compare <vaultA> <vaultB> --pass-a <passphrase> [--pass-b <passphrase>]
```

## Flags

| Flag | Description |
|------|-------------|
| `--pass-a` | Passphrase for vault A (required) |
| `--pass-b` | Passphrase for vault B (defaults to `--pass-a` if omitted) |
| `--keys-only` | Show only key names without counts or categories |

## Output

The command prints a summary of differences between the two vaults:

- **Only in A** — keys present in vault A but not in vault B
- **Only in B** — keys present in vault B but not in vault A
- **Changed** — keys present in both vaults but with different values
- **Identical** — keys present in both vaults with the same value

If the vaults are identical, the message `Vaults are identical.` is printed.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0`  | Vaults are identical |
| `1`  | Vaults differ |
| `2`  | Error (e.g. bad passphrase, missing file) |

This makes `envault compare` suitable for use in scripts and CI pipelines.

## Examples

### Compare two vaults with the same passphrase

```bash
envault compare staging.vault prod.vault --pass-a mysecret
```

### Compare two vaults with different passphrases

```bash
envault compare dev.vault prod.vault --pass-a devpass --pass-b prodpass
```

### Sample output

```
Only in A (1): DEBUG_MODE
Only in B (1): NEW_FEATURE_FLAG
Changed (2): DATABASE_URL, REDIS_URL
Identical (5): APP_NAME, LOG_LEVEL, PORT, REGION, TIMEOUT
```

## Notes

- Values are compared in plaintext after decryption; passphrases are never stored or logged.
- Use `envault diff` for a patch-style line-by-line view of a single vault's history.

# envault lintcheck

Run extended lint rules against a vault file to catch common key/value issues.

## Usage

```
envault lintcheck <vault-file> [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--level` | `-l` | `` | Filter results by severity: `warn` or `error` |

## Rules

| Rule | Level | Description |
|------|-------|-------------|
| `no-lowercase-key` | warn | Key contains lowercase letters; prefer `UPPER_SNAKE_CASE` |
| `no-empty-value` | warn | Value is empty |
| `no-spaces-in-key` | error | Key contains space characters |
| `no-special-chars-in-key` | error | Key contains characters outside `A-Z`, `0-9`, `_` |
| `no-numeric-prefix` | warn | Key starts with a digit |

## Examples

### Check a vault file

```bash
envault lintcheck .env.vault
```

Output:

```
LEVEL  KEY          RULE                  MESSAGE
error  BAD KEY      no-spaces-in-key      key contains spaces
warn   db_password  no-lowercase-key      key contains lowercase letters; prefer UPPER_SNAKE_CASE
warn   EMPTY_VAR    value is empty
```

### Show only errors

```bash
envault lintcheck --level error .env.vault
```

### Show only warnings

``` --level warn .env.vault
```

## Exit Codes

- `0` — No issues found, or only warnings present or more `error`-level findings detected

This makes `lintcheck` suitable for use in CI pipelines.

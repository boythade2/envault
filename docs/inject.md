# envault inject

Inject vault entries as environment variables into the current shell session or a subprocess.

## Usage

```
envault inject <vault-file> [flags] [-- command [args...]]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--overwrite` | `-o` | `false` | Overwrite existing environment variables |
| `--prefix` | `-p` | `""` | Only inject keys that start with this prefix |
| `--dry-run` | `-n` | `false` | Preview what would be injected without applying |

## Examples

### Inject all keys from a vault

```bash
envault inject myproject.vault
```

### Inject only keys with a specific prefix

```bash
envault inject myproject.vault --prefix APP_
```

### Preview what would be injected

```bash
envault inject myproject.vault --dry-run
```

### Run a command with injected environment

```bash
envault inject myproject.vault -- go run ./cmd/server
```

### Overwrite existing environment variables

```bash
envault inject myproject.vault --overwrite -- make test
```

## Output

After injection, a summary is printed:

```
Injected: 3, Overridden: 1, Skipped: 0
  + APP_HOST
  + APP_PORT
  + DB_NAME
  ~ LOG_LEVEL (overridden)
```

## Notes

- Without `--overwrite`, keys that already exist in the environment are skipped.
- When a command is provided after `--`, it is executed with the updated environment.
- The `--dry-run` flag is useful for auditing what a vault contains before applying it.

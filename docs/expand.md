# envault expand

Expand `${KEY}` references within vault entry values.

## Usage

```
envault expand <vault-file> [flags]
```

## Description

The `expand` command scans all values in a vault file for `${KEY}` placeholders
and resolves them by substituting the value of the referenced key from the **same
vault**.

After expansion the vault file is updated in place with the resolved values.

If a reference cannot be resolved within the vault, the command exits with an
error unless `--use-os` is supplied, in which case the OS environment is checked
as a fallback.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--use-os` | `false` | Fall back to OS environment variables for unresolved references |

## Examples

### Expand internal references

Given a vault that contains:

```
HOST=db.example.com
DB_URL=postgres://${HOST}/myapp
```

Running:

```bash
envault expand myproject.vault
```

Will update `DB_URL` to `postgres://db.example.com/myapp` and report:

```
expanded 1 reference(s):
  DB_URL: "postgres://${HOST}/myapp" -> "postgres://db.example.com/myapp"
```

### Use OS fallback

```bash
export REGION=us-east-1
envault expand myproject.vault --use-os
```

Any `${REGION}` placeholder that is not present in the vault will be resolved
from the shell environment.

## Notes

- Only `${KEY}` syntax is supported (not `$KEY`).
- References are resolved in a single pass; circular references are not detected.
- The vault file is written back only when at least one value changes.

# envault regex

Filter vault entries using a regular expression pattern.

## Usage

```
envault regex <vault-file> <pattern> [flags]
```

## Arguments

| Argument     | Description                          |
|--------------|--------------------------------------|
| `vault-file` | Path to the vault file               |
| `pattern`    | Regular expression to match against  |

## Flags

| Flag            | Default | Description                                  |
|-----------------|---------|----------------------------------------------|
| `--match-value` | false   | Also match the pattern against entry values  |

## Examples

### Filter by key prefix

```bash
envault regex .env.vault "^APP_"
```

Output:
```
KEY                            VALUE
APP_HOST                       localhost
APP_PORT                       8080
```

### Filter by value content

```bash
envault regex .env.vault "localhost" --match-value
```

Matches entries where either the key or value contains `localhost`.

### No matches

If no entries match the pattern, envault prints:

```
no entries matched
```

## Notes

- The pattern must be a valid Go regular expression.
- Without `--match-value`, only keys are tested against the pattern.
- Use anchors (`^`, `$`) for precise matching.

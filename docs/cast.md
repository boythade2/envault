# envault cast

Coerce vault entry values to a canonical type representation.

## Usage

```
envault cast <vault-file> <key>[,key,...] [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--type`, `-t` | `string` | Target type: `int`, `float`, `bool`, `string` |
| `--dry-run` | `false` | Preview changes without writing to disk |
| `--passphrase`, `-p` | `""` | Vault passphrase (or set `ENVAULT_PASSPHRASE`) |

## Supported Types

| Type | Example input | Normalised output |
|------|--------------|-------------------|
| `int` | `8080.0` | `8080` |
| `float` | `3` | `3` |
| `bool` | `True`, `1`, `yes` | `true` / `false` |
| `string` | any | unchanged |

## Examples

Cast a single key to integer:

```bash
envault cast prod.vault PORT --type int
```

Cast multiple keys to bool:

```bash
envault cast prod.vault DEBUG,VERBOSE --type bool
```

Preview changes without writing:

```bash
envault cast prod.vault PORT --type int --dry-run
```

## Output

```
KEY                      OLD              NEW              STATUS
PORT                     8080.0           8080             cast
DEBUG                    True             true             cast
NAME                     envault          envault          unchanged
```

## Notes

- If a value cannot be parsed as the target type an `error` status is shown and the entry is left unchanged.
- Multiple keys can be specified as a comma-separated list: `KEY1,KEY2,KEY3`.
- Use `ENVAULT_PASSPHRASE` environment variable to avoid passing the passphrase on the command line.

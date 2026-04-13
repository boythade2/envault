# envault archive

The `archive` command moves one or more keys from a vault into a separate archive file. Archived keys are removed from the active vault but their values are preserved for reference or auditing purposes.

## Usage

```
envault archive <vault-file> <key> [key...] [flags]
envault archive list <vault-file>
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--reason` | `""` | Optional reason for archiving the key(s) |
| `--dry-run` | `false` | Preview which keys would be archived without making changes |

## Examples

### Archive a single key

```bash
envault archive .env.vault OLD_API_KEY --reason "deprecated in v2"
```

### Archive multiple keys

```bash
envault archive .env.vault KEY_A KEY_B KEY_C
```

### Preview without writing

```bash
envault archive .env.vault LEGACY_KEY --dry-run
```

### List archived keys

```bash
envault archive list .env.vault
```

Output:

```
KEY                            ARCHIVED AT                    REASON
------------------------------------------------------------------------
OLD_API_KEY                    2024-06-01 12:00:00            deprecated in v2
```

## Archive file

Archived entries are stored in a hidden JSON file alongside the vault:

```
.env.vault.archive.json
```

The file is created with `0600` permissions and is not removed when new entries are archived — it accumulates all previously archived keys.

## Notes

- Archiving a key that does not exist in the vault returns an error.
- The archive file is **not** encrypted. Avoid archiving sensitive values in environments where the file may be exposed.
- Use `--dry-run` to safely inspect the impact before committing changes.

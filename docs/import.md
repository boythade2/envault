# `envault import`

Import environment variables from an existing `.env` or JSON file directly into your vault.

## Usage

```
envault import [file] [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--vault` | `-v` | `vault.json` | Path to the vault file |
| `--format` | `-f` | _(auto)_ | File format: `dotenv` or `json` |

## Format auto-detection

If `--format` is omitted, the format is inferred from the file extension:

- `.json` → `json`
- anything else (`.env`, `.txt`, etc.) → `dotenv`

## Examples

### Import a `.env` file

```bash
envault import .env
```

### Import a JSON file

```bash
envault import config/vars.json
```

### Specify vault path explicitly

```bash
envault import .env --vault /secrets/project.vault.json
```

### Force a format

```bash
envault import myfile --format dotenv
```

## Notes

- Lines beginning with `#` in dotenv files are treated as comments and skipped.
- Blank lines are ignored.
- Quoted values (e.g. `SECRET="abc"`) are unquoted automatically.
- Existing keys in the vault are **overwritten** by imported values.
- The vault file is created if it does not exist.

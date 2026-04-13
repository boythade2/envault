# envault group

Manage logical groupings of keys within a vault file. Groups allow you to
organise related environment variables (e.g. all database keys, all API keys)
without affecting their values or encryption.

## Subcommands

### `group add <vault> <group>`

Create a new empty group inside the vault.

```bash
envault group add .env.vault backend
```

### `group remove <vault> <group>`

Delete a group. Keys previously assigned to the group are not deleted from the
vault; only the group membership record is removed.

```bash
envault group remove .env.vault backend
```

### `group assign <vault> <group> <key>`

Add an existing vault key to a group.

```bash
envault group assign .env.vault backend DB_HOST
envault group assign .env.vault backend DB_PORT
```

### `group unassign <vault> <group> <key>`

Remove a key from a group without deleting the key from the vault.

```bash
envault group unassign .env.vault backend DB_HOST
```

### `group list <vault>`

Print all groups and their assigned keys in a table.

```bash
envault group list .env.vault
```

Example output:

```
GROUP    KEYS
backend  DB_HOST, DB_PORT
api      API_KEY
```

## Storage

Group metadata is stored in a sidecar file next to the vault:

```
.env.vault
.env.vault.groups.json   ← group membership (permissions 0600)
```

The groups file is plain JSON and is **not** encrypted. Do not store sensitive
information in group names or key names if the sidecar file may be exposed.

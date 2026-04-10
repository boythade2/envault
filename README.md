# envault

A lightweight CLI tool for managing and encrypting environment variable files across multiple projects.

---

## Installation

```bash
go install github.com/yourusername/envault@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envault.git && cd envault && go build -o envault .
```

---

## Usage

**Encrypt an `.env` file:**
```bash
envault encrypt --file .env --output .env.vault
```

**Decrypt a vault file:**
```bash
envault decrypt --file .env.vault --output .env
```

**Use with multiple projects:**
```bash
envault encrypt --file ./projects/api/.env --key my-secret-key
envault decrypt --file ./projects/api/.env.vault --key my-secret-key
```

Envault uses AES-256 encryption under the hood. Store your vault files safely in version control and share keys out-of-band.

---

## Commands

| Command    | Description                          |
|------------|--------------------------------------|
| `encrypt`  | Encrypt an environment variable file |
| `decrypt`  | Decrypt a vault file                 |
| `list`     | List all managed vault files         |
| `version`  | Print the current version            |

---

## License

[MIT](LICENSE) © 2024 yourusername
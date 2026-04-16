# envdiff

A CLI tool to compare `.env` files across environments and flag missing or mismatched keys.

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

## Usage

```bash
envdiff [flags] <base-file> <compare-file>
```

### Example

```bash
envdiff .env.example .env.production
```

**Sample output:**

```
MISSING KEYS in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED KEYS:
  - APP_ENV: "development" (base) vs "production" (compare)

✔ All other keys match.
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Only check for missing keys, skip value comparison |
| `--quiet` | Suppress output, exit code reflects result |
| `--json` | Output results in JSON format |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | No differences found |
| `1` | Differences detected |
| `2` | Error during execution |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)
# Contributing to AGEN

We love contributions! Here's how you can help.

## Development Setup

1. **Prerequisites**: Go 1.22+
2. **Clone**:
   ```bash
   git clone https://github.com/eshanized/agen.git
   ```
3. **Build**:
   ```bash
   make build
   ```

## Adding New Templates

Templates are located in `internal/templates/data`.

1. Add your agent markdown file to `internal/templates/data/agents/`.
2. Add frontmatter with metadata:
   ```yaml
   ---
   description: "My new agent"
   skills: docker, git
   ---
   ```
3. Rebuild the binary to embed the new files.

## Running Tests

```bash
make test
```

## Pull Requests

Please adhere to the Code of Conduct and ensure all tests pass before submitting a PR.

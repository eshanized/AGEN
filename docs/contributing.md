# Contributing to AGEN

We love contributions! This guide covers everything you need to know to contribute to AGEN.

## Development Setup

### Prerequisites

- **Go 1.22+** - [Install Go](https://go.dev/doc/install)
- **Git** - For version control
- **Make** (optional) - For development shortcuts

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/eshanized/agen.git
cd agen

# Install dependencies
go mod download

# Build the binary
go build -o agen ./cmd/agen

# Or use Make
make build
```

### Verify Installation

```bash
./agen --version
./agen list
```

---

## Project Structure

```
agen/
├── cmd/agen/              # Entry point
│   └── main.go
├── internal/              # Internal packages
│   ├── cli/              # CLI commands (18 files)
│   ├── ide/              # IDE adapters
│   ├── templates/        # Template engine
│   │   └── data/         # Embedded templates
│   ├── tui/              # Terminal UI
│   ├── ai/               # AI features
│   ├── verify/           # Verification
│   ├── plugin/           # Plugin manager
│   ├── team/             # Team config
│   ├── config/           # Global config
│   └── updater/          # Self-update
├── docs/                  # Documentation
├── templates/            # Example templates
└── Makefile
```

---

## Adding New Templates

Templates are the core content of AGEN.

### Adding an Agent

1. Create a new file in `internal/templates/data/agents/`:

```bash
touch internal/templates/data/agents/my-new-agent.md
```

2. Add YAML frontmatter and content:

```markdown
---
name: my-new-agent
description: Description of what this agent does and when to use it
tools: Read, Grep, Glob, Bash, Edit, Write
skills: clean-code, relevant-skill
---

# My New Agent

You are a specialized agent that...

## Your Expertise

...

## Guidelines

...
```

3. Rebuild the binary to embed the new template:

```bash
go build -o agen ./cmd/agen
```

### Adding a Skill

1. Create a skill directory:

```bash
mkdir -p internal/templates/data/skills/my-new-skill
touch internal/templates/data/skills/my-new-skill/SKILL.md
```

2. Add skill content:

```markdown
---
name: my-new-skill
description: Brief description of the skill
version: 1.0
---

# My New Skill Title

## Purpose

...

## Guidelines

...
```

3. Optionally add verification scripts:

```bash
mkdir internal/templates/data/skills/my-new-skill/scripts
touch internal/templates/data/skills/my-new-skill/scripts/my_checker.py
```

---

## Adding CLI Commands

### Command Structure

Commands are in `internal/cli/`. Each command follows this pattern:

```go
// internal/cli/mycommand.go
package cli

import (
    "fmt"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand <args>",
    Short: "Short description",
    Long:  `Longer description with examples.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runMyCommand,
}

func init() {
    // Register with root command
    rootCmd.AddCommand(myCmd)
    
    // Add flags
    myCmd.Flags().BoolP("force", "f", false, "force operation")
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    force, _ := cmd.Flags().GetBool("force")
    
    // Command logic here
    
    return nil
}
```

### Registering Commands

Commands are automatically registered via their `init()` functions. Just create the file and add the command to `rootCmd`.

---

## Adding IDE Adapters

### Adapter Interface

Implement the `Adapter` interface in `internal/ide/`:

```go
// internal/ide/myide.go
package ide

import (
    "github.com/eshanized/agen/internal/templates"
)

type MyIDEAdapter struct{}

func (a *MyIDEAdapter) Name() string {
    return "MyIDE"
}

func (a *MyIDEAdapter) Detect(projectPath string) bool {
    // Check for IDE-specific files
    // Return true if this IDE is detected
}

func (a *MyIDEAdapter) Install(tmpl *templates.Templates, opts InstallOptions) error {
    // Write templates in IDE-specific format
}

func (a *MyIDEAdapter) Update(tmpl *templates.Templates, opts UpdateOptions) (*UpdateChanges, error) {
    // Handle updates with conflict detection
}

func (a *MyIDEAdapter) GetRulesPath() string {
    return ".myide/rules"
}
```

### Registering Adapters

Add to `internal/ide/detector.go`:

```go
func init() {
    adapters["myide"] = &MyIDEAdapter{}
}
```

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/cli/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Or use Make
make test
```

---

## Code Style

### Formatting

```bash
# Format all code
go fmt ./...

# Or use Make
make fmt
```

### Linting

```bash
# Run linter (install golangci-lint first)
golangci-lint run

# Or use Make
make lint
```

### Guidelines

1. **Clear names**: Use descriptive variable and function names
2. **Comments**: Add comments for non-obvious logic
3. **Error handling**: Always handle errors, don't ignore them
4. **Small functions**: Keep functions focused and small
5. **Tests**: Add tests for new functionality

---

## Documentation

### Adding Documentation

Docs are in `docs/` and use [MkDocs](https://www.mkdocs.org/).

1. Create or edit markdown files in `docs/`
2. Update `mkdocs.yml` navigation if adding new pages
3. Preview locally:

```bash
pip install mkdocs-material pymdown-extensions
mkdocs serve
```

### Building Docs

```bash
mkdocs build
```

---

## Pull Request Process

### Before Submitting

1. **Fork** the repository
2. **Create a branch**: `git checkout -b feature/my-feature`
3. **Make changes** and commit with clear messages
4. **Run tests**: `go test ./...`
5. **Run linter**: `golangci-lint run`
6. **Update docs** if necessary

### Commit Messages

Use clear, descriptive commit messages:

```
feat: add new xyz command

- Implements the xyz functionality
- Adds tests for xyz
- Updates documentation
```

Prefixes:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation only
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks

### Submitting

1. Push your branch
2. Open a Pull Request
3. Fill out the PR template
4. Wait for review

---

## Release Process

Releases are automated via GoReleaser.

### Creating a Release

```bash
# Tag a new version
git tag v1.2.3
git push origin v1.2.3
```

GitHub Actions automatically:
1. Builds for all platforms
2. Creates GitHub release
3. Updates Homebrew tap
4. Updates Scoop bucket
5. Updates AUR package

---

## Getting Help

- **GitHub Issues**: Report bugs or request features
- **Discussions**: Ask questions or share ideas
- **Email**: eshanized@proton.me

---

## Code of Conduct

Please be respectful and constructive. We welcome contributors of all backgrounds and experience levels.

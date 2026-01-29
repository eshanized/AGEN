---
description: Documentation generation workflow for codebases.
---

# /document - Generate Documentation

$ARGUMENTS

---

## Purpose

Generate comprehensive documentation for a codebase.

---

## Task Flow

```
/document [scope]

1. DISCOVER
   └── Invoke explorer-agent
   └── Map codebase structure
   
2. ANALYZE
   └── Identify key components
   └── Extract API endpoints
   └── Find configuration files
   
3. GENERATE
   └── Invoke documentation-writer
   └── Create appropriate docs
   
4. VALIDATE
   └── Check links work
   └── Verify code examples run
```

---

## Documentation Types

### /document readme
Generate or update project README.

```markdown
# Project Name

Brief description.

## Features
- Feature 1
- Feature 2

## Quick Start
Installation and usage instructions.

## Documentation
Links to detailed docs.

## Contributing
How to contribute.

## License
License information.
```

### /document api
Generate API reference.

```markdown
# API Reference

## Authentication
How to authenticate.

## Endpoints

### GET /users
List all users.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|

**Response:**
```json
{
  "users": [...]
}
```
```

### /document architecture
Generate architecture documentation.

```markdown
# Architecture

## Overview
High-level system diagram.

## Components
Description of each component.

## Data Flow
How data moves through the system.

## Decisions
Key architectural decisions (ADRs).
```

### /document changelog
Generate or update changelog.

```markdown
# Changelog

## [Unreleased]

### Added
- New feature

### Changed
- Updated behavior

### Fixed
- Bug fix
```

---

## Output Locations

| Type | Location |
|------|----------|
| README | `./README.md` |
| API docs | `./docs/api.md` |
| Architecture | `./docs/architecture.md` |
| Changelog | `./CHANGELOG.md` |
| Contributing | `./CONTRIBUTING.md` |

---

## Usage Examples

```
/document                  # Generate all docs
/document readme           # Update README only
/document api              # Generate API reference
/document architecture     # Generate architecture docs
/document changelog        # Update changelog
/document --from-code      # Extract docs from code comments
```

---

## Best Practices

### Good Documentation
- Starts with "why", then "how"
- Includes working examples
- Stays up to date with code
- Has clear navigation

### Documentation Checklist
- [ ] README answers "what is this?"
- [ ] Quick start works in 5 minutes
- [ ] API reference is complete
- [ ] Examples are tested and working
- [ ] Links are not broken

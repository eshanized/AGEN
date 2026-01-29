# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Released]

## [2.0.0] - 2026-01-29

### Added

#### 🎯 8 New IDE Adapters (Total: 12 IDEs Supported)
- **Continue** - VS Code extension support (`.continue/`, `.continuerules`)
- **Cline** - VS Code extension support (`.clinerules`)
- **JetBrains** - IntelliJ, PyCharm, WebStorm support (`.idea/ai-assistant.xml`)
- **Neovim** - AI plugin support (`.nvim/ai-rules.md`)
- **Emacs** - gptel/ellama support (`.dir-locals.el`)
- **Aider** - CLI tool support (`.aider.conf.yml`)
- **Claude Code** - Anthropic's coding assistant (`CLAUDE.md`)
- **GitHub Copilot Workspace** - Custom instructions (`.github/copilot-instructions.md`)

#### 🤖 10 New Agent Templates (Total: 30+ Agents)
- `ai-ml-engineer` - ML pipelines, LLM integration, model fine-tuning
- `data-engineer` - ETL/ELT, data warehousing, pipeline design
- `accessibility-specialist` - WCAG compliance, ARIA, inclusive design
- `cloud-architect` - AWS/GCP/Azure, IaC, cost optimization
- `api-designer` - REST, GraphQL, OpenAPI, versioning
- `blockchain-developer` - Smart contracts, Web3, security
- `tech-lead` - Architecture decisions, code review, team enablement
- `embedded-systems-developer` - IoT, firmware, RTOS
- `refactoring-specialist` - Legacy modernization, technical debt
- `localization-specialist` - i18n/l10n, RTL support

#### 🛠 10 New Skills (Total: 46+ Skills)
- `llm-integration` - API patterns, prompt engineering, RAG
- `graphql-patterns` - Schema design, queries, mutations
- `kubernetes-patterns` - Resource selection, Helm charts
- `accessibility-patterns` - WCAG, ARIA, keyboard navigation
- `real-time-patterns` - WebSockets, SSE, CRDTs
- `monorepo-patterns` - Turborepo, Nx, workspace management
- `serverless-patterns` - Lambda, Edge Functions, cold starts
- `rust-patterns` - Ownership, error handling, async
- `observability-patterns` - Logging, metrics, distributed tracing
- `wasm-patterns` - Browser and server-side WebAssembly

#### 📋 8 New Workflow Commands (Total: 16 Workflows)
- `/audit` - Comprehensive project audits
- `/migrate` - Technology migration
- `/review` - Multi-agent code reviews
- `/onboard` - Developer onboarding documentation
- `/refactor` - Safe and systematic refactoring
- `/incident` - Production incident response
- `/prototype` - Rapid prototyping
- `/document` - Documentation generation

### Improved
- **87.2% Test Coverage** for IDE adapters
- Enhanced documentation with comprehensive IDE support guide
- Updated README with all 12 supported IDEs

## [1.0.0] - 2026-01-28

### Added
- Initial release of AGEN - AI Agent Template Manager.
- **Multi-IDE Support**: Auto-configuration for Antigravity, Cursor, Windsurf, and Zed.
- **Template Management**: Commands to list (`agen list`), search (`agen search`), and init (`agen init`) agent templates.
- **Health Checks**: `agen health` command to verify project configuration and dependencies.
- **Playground**: `agen playground` to create temporary test environments for agents.
- **Profiles**: `agen profile` to save and reuse agent configurations.
- **Verification**: Built-in verification runners for security, linting, and SEO.
- **Update System**: Smart conflict resolution when updating templates.
- **CI/CD**: GitHub Actions workflows for testing and building.
- **Installation**: Support for Homebrew, Scoop, AUR, and binary downloads.

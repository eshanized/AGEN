# Skill Catalog

AGEN includes **36 modular skills** that can be mixed and matched across agents. Skills are reusable capabilities that extend agent functionality.

## How Skills Work

Each skill is a folder containing:

- `SKILL.md` - Main skill definition with YAML frontmatter
- `scripts/` (optional) - Verification scripts for the skill

Skills are referenced by agents in their `skills:` frontmatter field.

---

## Development Skills

### clean-code

Pragmatic coding standards - concise, direct, no over-engineering. Core principles: SRP, DRY, KISS, YAGNI.

**Used by**: All agents

---

### code-review-checklist

Standards and checklists for reviewing code. Covers correctness, security, performance, and maintainability.

**Used by**: Backend Specialist, Frontend Specialist

---

### lint-and-validate

Code quality tools and linting configuration. Integrates with ESLint, Prettier, and language-specific linters.

**Scripts**: `lint_runner.py`, `type_coverage.py`

---

### tdd-workflow

Test-Driven Development workflow. Red-green-refactor cycle, test-first approach.

**Used by**: Test Engineer

---

### testing-patterns

Unit testing, integration testing, and E2E testing patterns. Test organization and mock strategies.

**Scripts**: `test_runner.py`

---

### systematic-debugging

Systematic approach to debugging. Root cause analysis, hypothesis testing, and isolation techniques.

**Used by**: Debugger

---

### parallel-agents

Patterns for working with multiple agents simultaneously. Coordination and context sharing.

**Used by**: Orchestrator

---

### intelligent-routing

Smart task routing to appropriate specialist agents based on context analysis.

**Used by**: Orchestrator

---

## Language & Framework Skills

### nextjs-react-expert

Next.js and React best practices. App Router, Server Components, hooks, and performance.

**Used by**: Frontend Specialist

---

### tailwind-patterns

Utility-first CSS patterns with Tailwind. Design tokens, responsive design, dark mode.

**Used by**: Frontend Specialist

---

### nodejs-best-practices

Node.js patterns and best practices. Error handling, async patterns, security.

**Used by**: Backend Specialist

---

### python-patterns

Pythonic code patterns. Type hints, virtual environments, packaging.

---

### bash-linux

Shell scripting and Linux administration. Automation, system management.

**Used by**: DevOps Engineer

---

### powershell-windows

Windows scripting with PowerShell. System administration and automation.

---

## Architecture & Design Skills

### api-patterns

REST and GraphQL API design patterns. Versioning, error handling, documentation.

**Scripts**: `api_validator.py`

**Used by**: Backend Specialist

---

### architecture

Software architecture patterns. Microservices, monoliths, event-driven, CQRS.

---

### database-design

Database schema design, normalization, indexing, and query optimization.

**Scripts**: `schema_validator.py`

**Used by**: Database Architect

---

### frontend-design

UI/UX principles and frontend architecture. Component design, state management.

**Scripts**: `ux_audit.py`, `accessibility_checker.py`

**Used by**: Frontend Specialist

---

### mobile-design

Mobile-specific UX patterns. Touch targets, navigation, platform conventions.

**Scripts**: `mobile_audit.py`

**Used by**: Mobile Developer

---

### web-design-guidelines

Modern web design principles. Typography, color, layout, accessibility.

**Used by**: Frontend Specialist

---

### game-development

Game loop patterns. Physics, rendering, input handling, game state.

**Used by**: Game Developer

---

## Infrastructure & DevOps Skills

### deployment-procedures

Release checklists and deployment strategies. Blue-green, canary, rollback.

**Used by**: DevOps Engineer

---

### server-management

Linux server operations. Configuration, monitoring, security hardening.

**Used by**: DevOps Engineer

---

### mcp-builder

Model Context Protocol configuration. MCP server setup and integration.

---

## Security Skills

### vulnerability-scanner

Security scanning tools and patterns. OWASP checks, dependency auditing.

**Scripts**: `security_scan.py`

**Used by**: Security Auditor, Penetration Tester

---

### red-team-tactics

Offensive security techniques. Penetration testing methodologies.

**Used by**: Penetration Tester

---

## Documentation & Content Skills

### documentation-templates

Standard documentation formats. README templates, API documentation, guides.

**Used by**: Documentation Writer

---

### plan-writing

Structured planning documents. PLAN.md format, milestones, deliverables.

**Used by**: Project Planner, Product Manager

---

### brainstorming

Ideation and brainstorming techniques. Mind mapping, problem decomposition.

**Used by**: Project Planner, Product Manager, Product Owner

---

## Specialized Skills

### seo-fundamentals

Search engine optimization. Meta tags, schema markup, content optimization.

**Scripts**: `seo_checker.py`

**Used by**: SEO Specialist

---

### geo-fundamentals

Geospatial data handling. Location services, mapping, coordinates.

**Scripts**: `geo_checker.py`

**Used by**: SEO Specialist

---

### i18n-localization

Internationalization and localization. Translation workflows, locale handling.

**Scripts**: `i18n_checker.py`

---

### performance-profiling

Performance analysis techniques. Profiling, benchmarking, optimization.

**Scripts**: `lighthouse_audit.py`

**Used by**: Performance Optimizer

---

### webapp-testing

Web application testing. Browser automation, Playwright, visual testing.

**Scripts**: `playwright_runner.py`

**Used by**: QA Automation Engineer

---

### app-builder

Application scaffolding patterns. Project structure, configuration.

---

### behavioral-modes

Agent behavioral modes. Planning, execution, verification phases.

---

## Usage

### Selecting Skills During Init

```bash
# Install specific skills
agen init --skills clean-code,testing-patterns,api-patterns

# Skills are automatically included with agents
agen init --agents backend-specialist  # Includes api-patterns, etc.
```

### Skill Files Location

After installation, skills are located at:

```
.agent/skills/<skill-name>/
├── SKILL.md          # Main skill definition
└── scripts/          # Optional verification scripts
    └── *.py
```

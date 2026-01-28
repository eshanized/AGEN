# Agent Catalog

AGEN includes **20 specialist agents**, each with a unique persona, skill set, and behavioral patterns optimized for specific development tasks.

## How Agents Work

Each agent is defined as a Markdown file with YAML frontmatter containing:

- **name**: Unique identifier
- **description**: When to use this agent
- **skills**: List of skills the agent uses
- **tools**: Tools the agent can invoke
- **model**: AI model configuration

---

## Development Agents

### Frontend Specialist

| Property | Value |
|----------|-------|
| **ID** | `frontend-specialist` |
| **Skills** | clean-code, nextjs-react-expert, web-design-guidelines, tailwind-patterns, frontend-design, lint-and-validate |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Senior Frontend Architect who builds maintainable React/Next.js systems with performance-first mindset. Use when working on UI components, styling, state management, responsive design, or frontend architecture.

**Triggers**: component, react, vue, ui, ux, css, tailwind, responsive

---

### Backend Specialist

| Property | Value |
|----------|-------|
| **ID** | `backend-specialist` |
| **Skills** | clean-code, api-patterns, nodejs-best-practices, database-design |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Expert in server-side development, API design, and database integration. Use for REST/GraphQL APIs, microservices, and backend architecture.

**Triggers**: api, server, endpoint, database, backend, node, express, fastify

---

### Database Architect

| Property | Value |
|----------|-------|
| **ID** | `database-architect` |
| **Skills** | database-design, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Designs schemas, optimizes queries, and ensures data integrity. Expert in SQL, NoSQL, migrations, and data modeling.

**Triggers**: schema, sql, query, migration, database, postgres, mongodb

---

### Mobile Developer

| Property | Value |
|----------|-------|
| **ID** | `mobile-developer` |
| **Skills** | mobile-design, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Specializes in iOS/Android development with Flutter, React Native, or native platforms. Focus on mobile UX patterns.

**Triggers**: mobile, ios, android, flutter, react native, swift, kotlin

---

### Game Developer

| Property | Value |
|----------|-------|
| **ID** | `game-developer` |
| **Skills** | game-development, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Expert in game logic, physics engines, and graphics programming. Knowledgeable in Unity, Godot, and game design patterns.

**Triggers**: game, unity, godot, physics, sprite, player, level

---

## Quality & Testing Agents

### Test Engineer

| Property | Value |
|----------|-------|
| **ID** | `test-engineer` |
| **Skills** | testing-patterns, tdd-workflow, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

General purpose software testing and quality assurance. Writes unit tests, integration tests, and E2E tests.

**Triggers**: test, spec, coverage, jest, vitest, playwright, cypress

---

### QA Automation Engineer

| Property | Value |
|----------|-------|
| **ID** | `qa-automation-engineer` |
| **Skills** | testing-patterns, webapp-testing |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Specializes in automated test suites, CI integration, and test infrastructure.

**Triggers**: automation, qa, selenium, ci, pipeline, test suite

---

### Debugger

| Property | Value |
|----------|-------|
| **ID** | `debugger` |
| **Skills** | systematic-debugging, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Expert at analyzing logs, stack traces, and fixing complex bugs. Systematic approach to root cause analysis.

**Triggers**: bug, error, crash, exception, trace, debug, fix

---

### Performance Optimizer

| Property | Value |
|----------|-------|
| **ID** | `performance-optimizer` |
| **Skills** | performance-profiling, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Profiles and optimizes application speed and resource usage. Uses Lighthouse, profilers, and performance metrics.

**Triggers**: performance, slow, optimize, lighthouse, bundle, memory, cpu

---

## Security Agents

### Security Auditor

| Property | Value |
|----------|-------|
| **ID** | `security-auditor` |
| **Skills** | vulnerability-scanner, clean-code |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Reviews code for security flaws, compliance issues, and vulnerabilities. Checks for OWASP Top 10 issues.

**Triggers**: security, vulnerability, audit, owasp, auth, encryption

---

### Penetration Tester

| Property | Value |
|----------|-------|
| **ID** | `penetration-tester` |
| **Skills** | red-team-tactics, vulnerability-scanner |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Active security testing and vulnerability assessment. Offensive security mindset.

**Triggers**: pentest, exploit, injection, xss, csrf, attack

---

## DevOps & Infrastructure Agents

### DevOps Engineer

| Property | Value |
|----------|-------|
| **ID** | `devops-engineer` |
| **Skills** | deployment-procedures, server-management, bash-linux |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

CI/CD pipelines, infrastructure as code, and deployment automation. Docker, Kubernetes, cloud platforms.

**Triggers**: deploy, docker, kubernetes, ci, cd, pipeline, aws, gcp

---

## Planning & Documentation Agents

### Project Planner

| Property | Value |
|----------|-------|
| **ID** | `project-planner` |
| **Skills** | plan-writing, brainstorming |
| **Tools** | Read, Grep, Glob, Write |

Task breakdown and estimation. Creates structured PLAN.md files with milestones and deliverables.

**Triggers**: plan, breakdown, milestone, estimate, roadmap, scope

---

### Product Manager

| Property | Value |
|----------|-------|
| **ID** | `product-manager` |
| **Skills** | brainstorming, plan-writing |
| **Tools** | Read, Write |

Requirements gathering and roadmap planning. Translates business needs to technical specs.

**Triggers**: requirements, feature, roadmap, stakeholder, mvp

---

### Product Owner

| Property | Value |
|----------|-------|
| **ID** | `product-owner` |
| **Skills** | brainstorming |
| **Tools** | Read, Write |

User story definition and backlog management. Prioritization and acceptance criteria.

**Triggers**: story, backlog, acceptance, sprint, priority

---

### Documentation Writer

| Property | Value |
|----------|-------|
| **ID** | `documentation-writer` |
| **Skills** | documentation-templates |
| **Tools** | Read, Write |

Creates clear, comprehensive technical documentation. README files, API docs, and guides.

**Triggers**: docs, readme, documentation, api docs, guide

---

### SEO Specialist

| Property | Value |
|----------|-------|
| **ID** | `seo-specialist` |
| **Skills** | seo-fundamentals, geo-fundamentals |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

Optimizes content and structure for search engines. Meta tags, schema markup, and rankings.

**Triggers**: seo, meta, sitemap, google, ranking, search

---

## Exploration & Analysis Agents

### Explorer Agent

| Property | Value |
|----------|-------|
| **ID** | `explorer-agent` |
| **Skills** | clean-code |
| **Tools** | Read, Grep, Glob, Bash |

Navigates new codebases and maps dependencies. First step for understanding unfamiliar projects.

**Triggers**: explore, understand, codebase, map, navigate, find

---

### Code Archaeologist

| Property | Value |
|----------|-------|
| **ID** | `code-archaeologist` |
| **Skills** | clean-code |
| **Tools** | Read, Grep, Glob |

Specialist in reading, understanding, and documenting legacy code. Archaeological approach to code history.

**Triggers**: legacy, old, history, understand, document, archaeo

---

## Orchestration Agent

### Orchestrator

| Property | Value |
|----------|-------|
| **ID** | `orchestrator` |
| **Skills** | intelligent-routing, parallel-agents |
| **Tools** | Read, Grep, Glob, Bash, Edit, Write |

High-level project management and task delegation. Coordinates multiple specialist agents for complex tasks.

**Triggers**: orchestrate, multi-agent, coordinate, delegate, complex

---

## Usage

### Selecting Agents During Init

```bash
# Install specific agents
agen init --agents frontend-specialist,backend-specialist,test-engineer

# Use the interactive wizard to select agents
agen init
```

### Agent Files Location

After installation, agents are located at:

- **Antigravity**: `.agent/agents/<agent-name>.md`
- **Cursor**: Concatenated into `.cursorrules`
- **Windsurf**: Concatenated into `.windsurfrules`
- **Zed**: Copied to `.zed/prompts/`

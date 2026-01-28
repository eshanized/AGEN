# Workflow Reference

AGEN includes **8 workflows** - predefined task automations that can be invoked with slash commands in supported IDEs.

## How Workflows Work

Workflows are Markdown files with YAML frontmatter containing:

- **description**: What the workflow does
- **Content**: Instructions executed when the workflow is invoked

Workflows are invoked using slash commands (e.g., `/orchestrate`, `/deploy`).

---

## Available Workflows

### /orchestrate

**Purpose**: Coordinate multiple agents for complex tasks.

Invokes the multi-agent orchestration system to:

1. Analyze task domains (security, frontend, backend, etc.)
2. Select 3+ appropriate specialist agents
3. Execute agents in parallel or sequence
4. Synthesize results into a unified report

**Usage**:
```
/orchestrate Build a full-stack e-commerce application with user authentication
```

**Key Features**:
- Minimum 3 agents required
- Two-phase approach: Planning → Implementation
- User approval checkpoint between phases
- Verification scripts run at the end

---

### /deploy

**Purpose**: Execute deployment procedures.

Guides through deployment process:

1. Pre-deployment checklist
2. Environment validation
3. Build and package
4. Deployment execution
5. Post-deployment verification

**Usage**:
```
/deploy to production with rollback plan
```

---

### /plan

**Purpose**: Create structured project plans.

Generates a `PLAN.md` file with:

- Problem statement
- Requirements analysis
- Technical approach
- Milestones and deliverables
- Risk assessment

**Usage**:
```
/plan Implement user authentication with OAuth
```

---

### /brainstorm

**Purpose**: Ideation and solution exploration.

Facilitates creative problem-solving:

1. Problem decomposition
2. Alternative approaches
3. Pros/cons analysis
4. Recommendation

**Usage**:
```
/brainstorm How should we structure our API versioning?
```

---

### /status

**Purpose**: Check project status and health.

Reports on:

- Installed agents and skills
- Configuration health
- Pending updates
- Recommendations

**Usage**:
```
/status
```

---

### /preview

**Purpose**: Preview changes before applying.

Shows what would change without making modifications:

- Files that would be created
- Files that would be modified
- Potential conflicts

**Usage**:
```
/preview adding security-auditor agent
```

---

### /enhance

**Purpose**: Improve existing code.

Applies best practices and improvements:

- Code quality enhancements
- Performance optimizations
- Accessibility improvements
- Documentation additions

**Usage**:
```
/enhance the UserProfile component
```

---

### /ui-ux-pro-max

**Purpose**: Advanced UI/UX design workflow.

Comprehensive design process:

1. Deep design thinking
2. Style commitment
3. Layout diversification
4. Animation and effects
5. Quality verification

**Key Constraints**:
- Purple is forbidden as primary color
- No default UI libraries without asking
- No standard/cliché designs
- Mandatory animations and visual depth

**Usage**:
```
/ui-ux-pro-max Design a landing page for a fintech startup
```

---

## Workflow Files Location

After installation, workflows are located at:

```
.agent/workflows/<workflow-name>.md
```

## Creating Custom Workflows

You can create custom workflows by adding Markdown files to `.agent/workflows/`:

```markdown
---
description: My custom workflow description
---

# My Custom Workflow

Instructions for the AI to follow when this workflow is invoked.

## Steps

1. First step
2. Second step
3. Third step
```

Custom workflows are automatically discovered and can be invoked using `/<filename>`.

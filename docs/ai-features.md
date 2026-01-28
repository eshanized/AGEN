# AI Features

AGEN includes AI-powered features that analyze your project and provide intelligent suggestions for agent selection and configuration.

## Overview

AI features include:

| Feature | Purpose |
|---------|---------|
| **Suggest** | Analyze project and recommend agents |
| **Explain** | Get detailed explanation of agents/skills |
| **Compose** | Create custom agents from description |

---

## Project Analysis

### Automatic Suggestions

```bash
agen ai suggest
```

AGEN analyzes your project to recommend appropriate agents.

### What Gets Analyzed

| Factor | Detection |
|--------|-----------|
| **Project Type** | Web, mobile, backend, CLI, library |
| **Languages** | JavaScript, TypeScript, Python, Go, etc. |
| **Frameworks** | React, Next.js, Express, Django, etc. |
| **Testing** | Has tests directory or test files |
| **Docker** | Has Dockerfile or docker-compose |
| **CI/CD** | Has .github/workflows, .gitlab-ci, etc. |
| **Dependencies** | Parses package.json, go.mod, etc. |

### Example Output

```
🔍 Analyzing project...

📊 Project Analysis
───────────────────
Type:       Web Application
Languages:  TypeScript, JavaScript
Frameworks: Next.js, React
Has Tests:  Yes
Has Docker: Yes
Has CI:     Yes

🎯 Recommended Agents
─────────────────────
Score  Agent                Reason
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
0.95   frontend-specialist  Next.js/React detected in dependencies
0.90   test-engineer        Test files found in project
0.85   devops-engineer      Docker and CI configuration present
0.80   security-auditor     Recommended for all web applications
0.70   performance-optimizer React project can benefit from optimization

💡 Run: agen init --agents frontend-specialist,test-engineer,devops-engineer
```

---

## Agent Explanation

### Get Detailed Info

```bash
agen ai explain frontend-specialist
```

Provides detailed information about an agent:

```
📦 Agent: frontend-specialist
─────────────────────────────

Senior Frontend Architect who builds maintainable React/Next.js 
systems with performance-first mindset.

🎯 Use Cases:
  • Building React/Next.js components or pages
  • Designing frontend architecture and state management
  • Implementing responsive UI or accessibility
  • Setting up styling (Tailwind, design systems)
  • Code reviewing frontend implementations

🧩 Skills Included:
  • clean-code
  • nextjs-react-expert
  • web-design-guidelines
  • tailwind-patterns
  • frontend-design
  • lint-and-validate

🔧 Tools Available:
  Read, Grep, Glob, Bash, Edit, Write

💡 Trigger Keywords:
  component, react, vue, ui, ux, css, tailwind, responsive
```

### Explain Skills

```bash
agen ai explain --skill clean-code
```

---

## Custom Agent Composition

### Create Custom Agents

```bash
agen ai compose my-reviewer \
  --description "Code reviewer focused on React and accessibility" \
  --base frontend-specialist,test-engineer
```

### What Gets Generated

AGEN creates a custom agent by:

1. Analyzing the description for key capabilities
2. Selecting relevant skills from base agents
3. Generating a custom system prompt
4. Creating a new agent markdown file

### Example Output

```
🎨 Composing custom agent: my-reviewer
────────────────────────────────────────

Based on analysis of:
  • Description: "Code reviewer focused on React and accessibility"
  • Base agents: frontend-specialist, test-engineer

📝 Generated Agent
──────────────────
Name: my-reviewer
Skills: clean-code, nextjs-react-expert, testing-patterns, 
        frontend-design, code-review-checklist

💾 Saved to: .agent/agents/my-reviewer.md

🔧 Customize by editing the generated file.
```

---

## Suggestion Scoring

Agents are scored 0.0 to 1.0 based on:

| Factor | Weight | Example |
|--------|--------|---------|
| **Framework Match** | High | React project → frontend-specialist |
| **Language Match** | Medium | Python → python-patterns skill |
| **Project Type** | Medium | API → backend-specialist |
| **Best Practice** | Low | All projects → security-auditor |

### Score Interpretation

| Score | Meaning |
|-------|---------|
| 0.90+ | Highly recommended, strong match |
| 0.70-0.89 | Recommended, good match |
| 0.50-0.69 | Consider, moderate match |
| Below 0.50 | Optional, weak match |

---

## Analysis Details

### Package.json Analysis

For JavaScript/TypeScript projects:

```json
{
  "dependencies": {
    "react": "18.x",      // → frontend-specialist
    "next": "14.x",       // → nextjs-react-expert skill
    "express": "4.x",     // → backend-specialist
    "jest": "29.x"        // → test-engineer
  }
}
```

### Go.mod Analysis

For Go projects:

```go
module myproject

require (
    github.com/gin-gonic/gin  // → backend-specialist
    github.com/lib/pq         // → database-architect
)
```

### File Pattern Detection

| Pattern | Detection |
|---------|-----------|
| `*.test.ts`, `*.spec.js` | Testing → test-engineer |
| `Dockerfile` | Docker → devops-engineer |
| `.github/workflows/` | CI/CD → devops-engineer |
| `*.sql`, `migrations/` | Database → database-architect |
| `*.md`, `docs/` | Documentation → documentation-writer |

---

## Best Practices

1. **Run on New Projects**: Use `agen ai suggest` when starting
2. **Review Suggestions**: Don't blindly accept all suggestions
3. **Iterate**: Re-run after major dependency changes
4. **Customize Compositions**: Edit generated custom agents
5. **Combine with Profiles**: Save good configurations as profiles

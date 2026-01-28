# Getting Started

Once you have installed AGEN, you are ready to start using it in your projects.

## Initialize a Project

Navigate to your project directory and run:

```bash
cd my-project
agen init
```

AGEN will:
1. **Detect your IDE**: It looks for `.cursorrules`, `.windsurfrules`, `.zed/`, or `.agent/` folders.
2. **Launch Wizard**: If no IDE is configured, it will ask you to choose one and select agents/skills.
3. **Install Templates**: It will copy the necessary configuration files to your project.

### Non-Interactive Mode

You can skip the wizard by passing flags:

```bash
# Force usage of Cursor IDE
agen init --ide cursor

# Install specific agents and skills
agen init --agents frontend,backend --skills docker
```

## List Available Templates

To see what agents and skills are available:

```bash
agen list
```

This will show:
- **Agents**: Specialist personalities (e.g., Frontend, Backend, DevOps)
- **Skills**: Reuseable capabilities (e.g., Docker, SQL, Git)
- **Workflows**: Common task automations

## Check Project Health

Run `agen health` to check if your project's agent configuration is healthy:

```bash
agen health
```

This checks for:
- Missing dependencies
- Broken links
- Syntax errors in configuration files

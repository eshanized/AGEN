# Advanced Usage

## Profiles

Profiles allow you to save a specific configuration (IDE + Agents + Skills) and reuse it later.

```bash
# Save current config as profile
agen profile save my-web-stack

# Apply profile to new project
agen init --profile my-web-stack
```

## Custom Templates

You can override built-in templates or add your own by creating a local `.agen-templates` directory in your home folder.

Structure:
```text
~/.agen-templates/
├── agents/
└── skills/
```

AGEN will look here first before using embedded templates.

## Conflict Resolution

When you run `agen update`, AGEN tries to preserve your local changes.
If a template has been modified both locally and in the new version, AGEN will prompt you to resolve the conflict or keep your local version.

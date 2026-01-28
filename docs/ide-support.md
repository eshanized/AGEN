# IDE Support

AGEN supports multiple IDEs by adapting its internal template format to the specific requirements of each environment.

## Antigravity (Claude Code)

This is the native format for AGEN.

- **Detection**: Checks for `.agent/` directory or `GEMINI.md`.
- **Structure**:
  ```text
  .agent/
  ├── agents/       # Agent personas (.md)
  ├── skills/       # Reusable skills (folders with SKILL.md)
  ├── workflows/    # Workflows (.md)
  └── rules/        # Global rules
  ```
- **Features**: Full support for all AGEN features.

## Cursor

Cursor uses a single rules file to define AI behavior.

- **Detection**: Checks for `.cursorrules`.
- **Format**: Concatenates all selected agents and skills into a single `.cursorrules` file.
- **Sections**:
  - Global Project Context
  - Selected Agents
  - Skills & Tool definitions

## Windsurf

Windsurf follows a similar pattern to Cursor but with its own rules file.

- **Detection**: Checks for `.windsurfrules`.
- **Format**: Single `.windsurfrules` file containing all context.

## Zed

Zed uses a folder-based approach for prompts and context.

- **Detection**: Checks for `.zed/` directory.
- **Structure**:
  ```text
  .zed/
  ├── settings.json   # Configuration
  └── prompts/        # Context files
  ```

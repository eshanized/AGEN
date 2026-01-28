# Verification System

AGEN includes a built-in verification system that performs automated checks on your project for security, code quality, UX, and SEO.

## Overview

The verification runner (`agen verify`) executes multiple check categories:

| Check | Purpose |
|-------|---------|
| **Security** | Scans for hardcoded secrets, vulnerabilities |
| **Lint** | Runs code quality checks |
| **UX** | Audits accessibility and usability |
| **SEO** | Checks search engine optimization |

---

## Running Verification

### Basic Usage

```bash
# Run all verification checks
agen verify

# Run with verbose output
agen verify --verbose
```

### Individual Checks

The verification system runs all checks by default. Results are displayed with:

- ✅ **Pass**: Check completed successfully
- ⚠️ **Warning**: Non-critical issues found
- ❌ **Critical**: Issues that must be fixed

---

## Security Scanning

Scans your codebase for security issues:

### What It Checks

| Check | Description |
|-------|-------------|
| **Hardcoded Secrets** | API keys, passwords, tokens in code |
| **Environment Files** | `.env` files not in `.gitignore` |
| **Security Anti-patterns** | Common vulnerability patterns |
| **Dependency Risks** | Known vulnerable packages |

### Patterns Detected

```
# API Keys
api_key = "sk-1234..."
API_KEY: "1234..."

# Passwords
password = "secret123"
DB_PASSWORD=mypass

# Tokens
token = "ghp_..."
jwt_secret = "..."
```

---

## Lint Checking

Runs language-specific linters on your codebase.

### JavaScript/TypeScript Projects

Attempts to run:
1. `npm run lint` if script exists
2. Falls back to `npx eslint .`
3. Manual checks if no linter available

### Go Projects

Attempts to run:
1. `golangci-lint run`
2. Falls back to `go vet ./...`

### Other Languages

Basic file checks for common issues:
- Trailing whitespace
- Mixed line endings
- Syntax errors in config files

---

## UX Auditing

Checks HTML/JSX files for usability issues:

### What It Checks

| Check | Description |
|-------|-------------|
| **Alt Text** | Images missing `alt` attributes |
| **ARIA Labels** | Interactive elements without labels |
| **Touch Targets** | Buttons/links too small (< 44px) |
| **Contrast** | Basic color contrast issues |

### Example Issues

```html
<!-- Missing alt text -->
<img src="photo.jpg">  ❌
<img src="photo.jpg" alt="User profile">  ✅

<!-- Missing aria-label -->
<button><icon/></button>  ❌
<button aria-label="Close dialog"><icon/></button>  ✅
```

---

## SEO Checking

Checks HTML files for search engine optimization:

### What It Checks

| Check | Description |
|-------|-------------|
| **Title Tag** | Page has `<title>` element |
| **Meta Description** | Has meta description |
| **OG Tags** | Open Graph tags for social sharing |
| **Canonical URL** | Has canonical link |
| **H1 Usage** | Proper heading hierarchy |

### Example Issues

```html
<!-- Missing title -->
<head></head>  ❌
<head><title>My App</title></head>  ✅

<!-- Missing meta description -->
<head>
  <meta name="description" content="...">  ✅
</head>
```

---

## Skill Verification Scripts

Beyond `agen verify`, individual skills include Python verification scripts:

### Available Scripts

| Skill | Script | Purpose |
|-------|--------|---------|
| frontend-design | `ux_audit.py` | UX analysis |
| frontend-design | `accessibility_checker.py` | A11y checks |
| api-patterns | `api_validator.py` | API validation |
| mobile-design | `mobile_audit.py` | Mobile UX |
| database-design | `schema_validator.py` | Schema validation |
| vulnerability-scanner | `security_scan.py` | Security scan |
| seo-fundamentals | `seo_checker.py` | SEO checks |
| geo-fundamentals | `geo_checker.py` | GEO checks |
| performance-profiling | `lighthouse_audit.py` | Lighthouse |
| testing-patterns | `test_runner.py` | Test execution |
| webapp-testing | `playwright_runner.py` | Browser tests |
| lint-and-validate | `lint_runner.py` | Linting |
| lint-and-validate | `type_coverage.py` | Type coverage |
| i18n-localization | `i18n_checker.py` | i18n checks |

### Running Skill Scripts

```bash
# From project root
python .agent/skills/vulnerability-scanner/scripts/security_scan.py .
python .agent/skills/lint-and-validate/scripts/lint_runner.py .
```

---

## Verification Results

Results are displayed with severity levels:

```
🔍 Running Security Scan...

❌ CRITICAL (2 items)
  • src/config.js:15 - Hardcoded API key detected
  • .env:3 - Database password exposed

⚠️ WARNING (1 item)
  • package.json - Vulnerable dependency: lodash < 4.17.21

✅ PASSED (5 items)
  • No exposed secrets in environment
  • .gitignore properly configured
  • No SQL injection patterns
  • No XSS vulnerabilities
  • Dependencies mostly up-to-date
```

---

## Best Practices

1. **Run Before Commit**: Add `agen verify` to pre-commit hooks
2. **CI Integration**: Include verification in your CI pipeline
3. **Fix Critical First**: Address critical issues before warnings
4. **Regular Audits**: Run verification regularly, not just at release

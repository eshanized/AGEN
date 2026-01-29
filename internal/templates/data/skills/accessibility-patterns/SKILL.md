---
name: accessibility-patterns
description: WCAG 2.2 compliance, ARIA patterns, keyboard navigation, and inclusive design patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Accessibility Patterns

> Principles for building accessible web applications.
> **Design for everyone, not just the average user.**

---

## 1. WCAG 2.2 Quick Reference

### Compliance Levels

| Level | What It Means |
|-------|---------------|
| **A** | Minimum accessibility |
| **AA** | Standard (legal requirement) |
| **AAA** | Enhanced accessibility |

### POUR Principles

| Principle | Focus |
|-----------|-------|
| **Perceivable** | All users can perceive content |
| **Operable** | All users can interact |
| **Understandable** | Content is clear |
| **Robust** | Works with assistive tech |

---

## 2. Color & Contrast

### Minimum Ratios

| Content | Ratio |
|---------|-------|
| **Normal text** | 4.5:1 |
| **Large text (18px+)** | 3:1 |
| **UI components** | 3:1 |

### Testing

```bash
# Tools
- Chrome DevTools → Accessibility panel
- axe DevTools browser extension
- WebAIM Contrast Checker
```

---

## 3. Keyboard Navigation

### Requirements

| Requirement | Implementation |
|-------------|----------------|
| **Tab-focusable** | All interactive elements |
| **Logical order** | Follows visual order |
| **No traps** | Can always escape |
| **Visible focus** | Clear focus indicator |
| **Skip links** | Skip to main content |

### Focus Management

```typescript
// After modal opens
modal.querySelector('[data-autofocus]').focus();

// After modal closes
triggerButton.focus();
```

### Focus Styles

```css
/* Visible focus for keyboard users */
:focus-visible {
  outline: 2px solid #0066ff;
  outline-offset: 2px;
}

/* Remove default for mouse users */
:focus:not(:focus-visible) {
  outline: none;
}
```

---

## 4. Semantic HTML

### Correct Elements

| Need | Use | Not |
|------|-----|-----|
| **Clickable** | `<button>` | `<div onclick>` |
| **Navigation** | `<nav>` | `<div class="nav">` |
| **Main content** | `<main>` | `<div id="main">` |
| **Article** | `<article>` | `<div class="article">` |
| **List** | `<ul>/<ol>` | `<div>` with CSS |

### Heading Hierarchy

```html
<!-- Correct -->
<h1>Page Title</h1>
  <h2>Section</h2>
    <h3>Subsection</h3>
  <h2>Another Section</h2>

<!-- Wrong: Skipping levels -->
<h1>Page Title</h1>
  <h3>Section</h3>  <!-- Missing h2 -->
```

---

## 5. ARIA Patterns

### When to Use ARIA

```
Native HTML can do it? → Don't use ARIA
Custom component? → Add ARIA
Dynamic content? → Use aria-live
```

### Common Patterns

| Component | ARIA |
|-----------|------|
| **Modal** | `role="dialog"`, `aria-modal="true"` |
| **Tabs** | `role="tablist/tab/tabpanel"` |
| **Alert** | `role="alert"` or `aria-live="assertive"` |
| **Loading** | `aria-busy="true"`, `aria-live="polite"` |
| **Toggle** | `aria-pressed="true/false"` |

### Live Regions

```html
<!-- Polite: Announce when convenient -->
<div aria-live="polite">
  Loading complete. 5 results found.
</div>

<!-- Assertive: Announce immediately -->
<div aria-live="assertive" role="alert">
  Error: Form submission failed.
</div>
```

---

## 6. Forms

### Labels

```html
<!-- Explicit label -->
<label for="email">Email</label>
<input type="email" id="email" name="email">

<!-- Group related fields -->
<fieldset>
  <legend>Shipping Address</legend>
  <!-- fields -->
</fieldset>
```

### Error Messages

```html
<label for="email">Email</label>
<input 
  type="email" 
  id="email" 
  aria-describedby="email-error"
  aria-invalid="true"
>
<span id="email-error" role="alert">
  Please enter a valid email address.
</span>
```

---

## 7. Images & Media

### Alt Text

| Image Type | Alt Text |
|------------|----------|
| **Informative** | Describe content |
| **Decorative** | `alt=""` (empty) |
| **Complex** | Brief alt + longer description |
| **Link/Button** | Describe action |

### Video/Audio

- Captions for video
- Transcripts for audio
- Audio description for visual content

---

## 8. Motion & Preferences

### Reduced Motion

```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## ✅ Checklist

- [ ] **Color contrast meets 4.5:1?**
- [ ] **All interactive elements keyboard accessible?**
- [ ] **Focus visible and logical?**
- [ ] **Semantic HTML used?**
- [ ] **Images have alt text?**
- [ ] **Forms have labels?**
- [ ] **Errors announced?**
- [ ] **Reduced motion respected?**
- [ ] **Tested with screen reader?**

---

> **Remember:** Accessibility benefits everyone. Curb cuts help wheelchairs, strollers, and rolling luggage alike.

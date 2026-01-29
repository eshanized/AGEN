---
description: Rapid prototyping workflow for MVPs and proof of concepts.
---

# /prototype - Rapid Prototyping

$ARGUMENTS

---

## Purpose

Quickly build a working prototype to validate ideas before full implementation.

---

## 🔴 CRITICAL RULES

1. **Speed over perfection** - It's a prototype, not production
2. **Validate hypothesis** - Build the minimum to test your idea
3. **Mock what you can** - Don't build what you don't need yet
4. **Document limitations** - Be clear about what's missing

---

## Task Flow

```
/prototype [idea]

1. SCOPE
   └── What are we validating?
   └── Minimum features needed
   
2. SELECT STACK
   └── Fastest path to working demo
   └── Use familiar tools
   
3. BUILD
   └── Core functionality first
   └── Mock secondary features
   
4. DEPLOY
   └── Quick deployment (Vercel, Netlify)
   └── Share with stakeholders
   
5. DOCUMENT
   └── What works
   └── What's mocked
   └── What's needed for production
```

---

## Prototyping Stack

### Frontend
| Type | Tool | Why |
|------|------|-----|
| **Web** | Next.js + shadcn | Fast, good defaults |
| **Mobile** | Expo | Quick React Native setup |
| **AI Chat** | Vercel AI SDK | Streaming, providers |

### Backend
| Need | Tool | Why |
|------|------|-----|
| **API** | Next.js API routes | Co-located |
| **Auth** | NextAuth/Clerk | Drop-in |
| **Database** | SQLite / Supabase | Quick setup |

### Mocking
| What | How |
|------|-----|
| **API responses** | JSON files, MSW |
| **Database** | In-memory arrays |
| **Auth** | Hardcoded user |
| **Payments** | Success always |

---

## Output Format

```markdown
# Prototype: [Name]

## Hypothesis
[What we're testing]

## Demo
[URL or screenshot]

---

## Features

### ✅ Working
- Feature A (fully functional)
- Feature B (fully functional)

### 🟡 Mocked
- Payment processing (always succeeds)
- Email sending (logged to console)

### ❌ Not Implemented
- User settings
- Admin panel

---

## Tech Stack
- Frontend: Next.js 14
- Styling: Tailwind + shadcn/ui
- Database: SQLite (file-based)
- Deploy: Vercel

---

## Path to Production

To make this production-ready:

1. **Required**
   - [ ] Real authentication
   - [ ] Production database
   - [ ] Payment integration

2. **Recommended**
   - [ ] Error monitoring
   - [ ] Analytics
   - [ ] Tests

3. **Estimated Effort**
   - 2-3 weeks for production MVP
```

---

## Usage Examples

```
/prototype "AI writing assistant"
/prototype "e-commerce checkout flow"
/prototype --mobile "location-based app"
```

---

## Key Principles

### What to Skip in Prototypes
- Perfect styling
- Edge case handling
- Comprehensive testing
- Performance optimization
- Full security audit

### What to Include
- Core user flow
- Key interactions
- Enough polish to demo
- Documentation of limitations

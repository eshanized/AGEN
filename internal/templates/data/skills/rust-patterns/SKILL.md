---
name: rust-patterns
description: Rust idioms, ownership, lifetimes, error handling, and async patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Rust Patterns

> Principles for writing idiomatic, safe Rust code.
> **Ownership is not a burden—it's a superpower.**

---

## 1. Ownership & Borrowing

### Core Rules

| Rule | Meaning |
|------|---------|
| **One owner** | Each value has exactly one owner |
| **Borrow or move** | Pass by reference or transfer ownership |
| **No aliasing + mutation** | One mutable OR many immutable refs |

### Common Patterns

```rust
// ✅ Borrow when you just need to read
fn print_name(name: &str) {
    println!("{}", name);
}

// ✅ Take ownership when you need to store
fn store_name(name: String) -> Person {
    Person { name }
}

// ✅ Return owned data from functions
fn create_greeting(name: &str) -> String {
    format!("Hello, {}!", name)
}
```

---

## 2. Error Handling

### Result & Option

```rust
// Prefer Result for recoverable errors
fn parse_config(path: &str) -> Result<Config, ConfigError> {
    let contents = fs::read_to_string(path)?;
    let config = serde_json::from_str(&contents)?;
    Ok(config)
}

// Option for optional values
fn find_user(id: u64) -> Option<User> {
    users.iter().find(|u| u.id == id).cloned()
}
```

### The ? Operator

```rust
// Propagate errors elegantly
fn process() -> Result<Output, Error> {
    let data = fetch_data()?;    // Returns early if Err
    let parsed = parse(data)?;   // Returns early if Err
    Ok(transform(parsed))
}
```

### Custom Errors

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),
    
    #[error("Not found: {0}")]
    NotFound(String),
    
    #[error("Invalid input: {0}")]
    Validation(String),
}
```

---

## 3. Structs & Enums

### Builder Pattern

```rust
#[derive(Default)]
pub struct ServerBuilder {
    port: Option<u16>,
    host: Option<String>,
}

impl ServerBuilder {
    pub fn port(mut self, port: u16) -> Self {
        self.port = Some(port);
        self
    }
    
    pub fn build(self) -> Result<Server, Error> {
        Ok(Server {
            port: self.port.unwrap_or(8080),
            host: self.host.unwrap_or_else(|| "127.0.0.1".into()),
        })
    }
}
```

### Newtype Pattern

```rust
// Type safety without runtime cost
pub struct UserId(u64);
pub struct OrderId(u64);

// Can't accidentally mix up IDs
fn get_order(user_id: UserId, order_id: OrderId) -> Order { ... }
```

---

## 4. Traits

### Common Traits to Implement

| Trait | When |
|-------|------|
| `Debug` | Always (derive) |
| `Clone` | Value needs copying |
| `Default` | Has sensible default |
| `Display` | User-facing output |
| `From/Into` | Type conversions |
| `Serialize/Deserialize` | JSON/data handling |

### Trait Objects vs Generics

```rust
// Generics: Static dispatch, faster
fn process<T: Handler>(handler: T) { ... }

// Trait objects: Dynamic dispatch, flexibility
fn process(handler: Box<dyn Handler>) { ... }
```

---

## 5. Async Rust

### Runtime Selection

| Runtime | Best For |
|---------|----------|
| **Tokio** | General purpose, most popular |
| **async-std** | std-like API |
| **smol** | Minimal, embeddable |

### Common Patterns

```rust
// Spawn concurrent tasks
let (result1, result2) = tokio::join!(
    fetch_users(),
    fetch_orders()
);

// Select first to complete
tokio::select! {
    result = operation() => handle(result),
    _ = tokio::time::sleep(timeout) => handle_timeout(),
}
```

### Avoid Blocking

```rust
// ❌ Don't block in async
async fn bad() {
    std::thread::sleep(Duration::from_secs(1)); // Blocks executor!
}

// ✅ Use async sleep
async fn good() {
    tokio::time::sleep(Duration::from_secs(1)).await;
}
```

---

## 6. Performance

### Zero-Cost Abstractions

```rust
// Iterators compile to optimal loops
let sum: u64 = numbers
    .iter()
    .filter(|n| **n > 0)
    .map(|n| n * 2)
    .sum();
```

### Allocation Tips

| Do | Don't |
|-----|-------|
| `&str` for read-only | `String` everywhere |
| `Vec::with_capacity` | Many `push` calls |
| `Cow<str>` | Clone just in case |

---

## ✅ Checklist

- [ ] **Ownership clear?** (each value has one owner)
- [ ] **Errors handled?** (Result, no unwrap in prod)
- [ ] **Lifetimes explicit?** (when needed)
- [ ] **Traits derived?** (Debug, Clone where appropriate)
- [ ] **No blocking in async?** (use async equivalents)
- [ ] **Allocations minimized?** (borrowed where possible)

---

> **Remember:** The compiler is your friend. If it compiles, it probably works. Trust the borrow checker.

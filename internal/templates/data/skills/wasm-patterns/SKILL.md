---
name: wasm-patterns
description: WebAssembly integration, WASI, and browser/server WASM patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# WebAssembly Patterns

> Principles for WebAssembly development.
> **Performance-critical code. Portable binaries. Language interop.**

---

## 1. When to Use WASM

### Good Use Cases

| Use Case | Why |
|----------|-----|
| **CPU-intensive** | Image/video processing |
| **Existing code** | Port C/C++/Rust to web |
| **Security** | Sandboxed execution |
| **Plugins** | Isolated user code |

### Not Ideal For

| Case | Why Not |
|------|---------|
| **DOM manipulation** | JS is faster for this |
| **Simple logic** | Overhead not worth it |
| **Network I/O** | Still needs JS interop |

---

## 2. Language Selection

| Language | Best For |
|----------|----------|
| **Rust** | wasm-bindgen, wasm-pack, best tooling |
| **AssemblyScript** | TypeScript developers |
| **C/C++** | Emscripten, existing codebases |
| **Go** | TinyGo for small binaries |

---

## 3. Browser WASM

### Setup with Rust

```bash
# Install toolchain
rustup target add wasm32-unknown-unknown
cargo install wasm-pack

# Build
wasm-pack build --target web
```

### JavaScript Integration

```javascript
import init, { process_image } from './pkg/my_wasm.js';

async function main() {
  await init();  // Initialize WASM module
  
  const result = process_image(imageData);
  console.log(result);
}
```

### wasm-bindgen Example

```rust
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[wasm_bindgen]
pub struct Counter {
    count: i32,
}

#[wasm_bindgen]
impl Counter {
    pub fn new() -> Counter {
        Counter { count: 0 }
    }
    
    pub fn increment(&mut self) {
        self.count += 1;
    }
}
```

---

## 4. WASI (Server-Side)

### What is WASI?

WebAssembly System Interface - run WASM outside browser with system access.

### Runtimes

| Runtime | Features |
|---------|----------|
| **Wasmtime** | Fast, Bytecode Alliance |
| **Wasmer** | Universal, many languages |
| **WasmEdge** | Edge/cloud, CNCF project |

### Build for WASI

```bash
# Rust target
rustup target add wasm32-wasi
cargo build --target wasm32-wasi

# Run
wasmtime target/wasm32-wasi/debug/my-app.wasm
```

---

## 5. Performance Optimization

### Bundle Size

| Technique | Effect |
|-----------|--------|
| `wasm-opt -O3` | Optimize binary |
| `lto = true` | Link-time optimization |
| Avoid `std` | Use `no_std` when possible |
| Remove panic strings | `panic = "abort"` |

### Cargo.toml

```toml
[profile.release]
lto = true
opt-level = "s"  # Size optimization
panic = "abort"
codegen-units = 1
```

---

## 6. Memory Management

### Linear Memory

```javascript
// Access WASM memory from JS
const memory = wasmInstance.exports.memory;
const buffer = new Uint8Array(memory.buffer);

// Pass data to WASM
const ptr = wasmInstance.exports.alloc(data.length);
buffer.set(data, ptr);
wasmInstance.exports.process(ptr, data.length);
```

### Best Practices

| Do | Don't |
|----|-------|
| Pre-allocate | Many small allocations |
| Reuse memory | Create/destroy often |
| Copy in bulk | Byte-by-byte transfer |

---

## 7. Component Model

### Future of WASM

```wit
// WIT (WASM Interface Types)
package example:calc;

interface operations {
  add: func(a: s32, b: s32) -> s32;
}

world calculator {
  export operations;
}
```

---

## ✅ Checklist

- [ ] **Use case justified?** (CPU-intensive, port code)
- [ ] **Bundle size optimized?** (wasm-opt, LTO)
- [ ] **Memory managed?** (no leaks, efficient transfer)
- [ ] **Error handling works?** (across boundary)
- [ ] **Async considered?** (Web Workers if blocking)
- [ ] **Fallback exists?** (JS implementation)

---

> **Remember:** WASM is a tool, not a goal. Use it when JavaScript isn't fast enough, not just because it's cool.

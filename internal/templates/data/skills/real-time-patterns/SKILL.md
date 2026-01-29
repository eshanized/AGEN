---
name: real-time-patterns
description: WebSocket, Server-Sent Events, real-time sync, and live collaboration patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Real-Time Patterns

> Principles for building real-time web applications.
> **Choose simplicity over complexity. SSE before WebSocket.**

---

## 1. Technology Selection

### Decision Matrix

| Need | Technology |
|------|------------|
| **Server → Client only** | Server-Sent Events (SSE) |
| **Bidirectional** | WebSocket |
| **Conflict resolution** | CRDTs (Yjs, Automerge) |
| **Presence/Typing** | WebSocket + ephemeral state |

### Comparison

| Feature | SSE | WebSocket |
|---------|-----|-----------|
| **Direction** | One-way | Bidirectional |
| **Protocol** | HTTP | WS |
| **Reconnection** | Built-in | Manual |
| **Complexity** | Low | Medium |
| **Binary data** | No (needs encoding) | Yes |

---

## 2. Server-Sent Events

### Server (Node.js)

```typescript
app.get('/events', (req, res) => {
  res.setHeader('Content-Type', 'text/event-stream');
  res.setHeader('Cache-Control', 'no-cache');
  res.setHeader('Connection', 'keep-alive');

  // Send event
  const sendEvent = (data) => {
    res.write(`data: ${JSON.stringify(data)}\n\n`);
  };

  // Heartbeat
  const heartbeat = setInterval(() => {
    res.write(': heartbeat\n\n');
  }, 30000);

  req.on('close', () => {
    clearInterval(heartbeat);
  });
});
```

### Client

```typescript
const eventSource = new EventSource('/events');

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  // Handle update
};

eventSource.onerror = () => {
  // Auto-reconnects
};
```

---

## 3. WebSocket Patterns

### Connection Management

```typescript
class WebSocketClient {
  private ws: WebSocket;
  private reconnectAttempts = 0;

  connect() {
    this.ws = new WebSocket(url);
    
    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
    };
    
    this.ws.onclose = () => {
      this.scheduleReconnect();
    };
  }

  private scheduleReconnect() {
    const delay = Math.min(1000 * 2 ** this.reconnectAttempts, 30000);
    setTimeout(() => this.connect(), delay);
    this.reconnectAttempts++;
  }
}
```

### Message Protocol

```typescript
interface Message {
  type: string;
  payload: unknown;
  id?: string;  // For acknowledgment
}

// Types
type MessageType = 
  | 'subscribe' 
  | 'unsubscribe' 
  | 'update' 
  | 'ack' 
  | 'error';
```

---

## 4. Presence & Typing Indicators

### Ephemeral State

```typescript
// Typing indicator with debounce
let typingTimeout;

onInput(() => {
  send({ type: 'typing', userId });
  
  clearTimeout(typingTimeout);
  typingTimeout = setTimeout(() => {
    send({ type: 'stop-typing', userId });
  }, 1000);
});
```

### Presence

```typescript
// Track online users
const presence = new Map();

ws.on('message', (msg) => {
  if (msg.type === 'presence') {
    presence.set(msg.userId, {
      lastSeen: Date.now(),
      status: msg.status
    });
  }
});

// Clean up stale presence
setInterval(() => {
  const now = Date.now();
  for (const [id, data] of presence) {
    if (now - data.lastSeen > 30000) {
      presence.delete(id);
    }
  }
}, 10000);
```

---

## 5. CRDTs for Collaboration

### When to Use

| Scenario | Solution |
|----------|----------|
| **Simple updates** | Last-write-wins |
| **Text editing** | Yjs, Automerge |
| **Counters** | G-Counter, PN-Counter |
| **Sets** | OR-Set, LWW-Set |

### Yjs Example

```typescript
import * as Y from 'yjs';
import { WebsocketProvider } from 'y-websocket';

const doc = new Y.Doc();
const provider = new WebsocketProvider(
  'wss://server.com',
  'room-name',
  doc
);

const text = doc.getText('content');
text.observe(() => {
  // Update UI
});
```

---

## 6. Scaling WebSockets

### Strategies

| Strategy | Use Case |
|----------|----------|
| **Pub/Sub (Redis)** | Multi-server broadcast |
| **Sticky Sessions** | Simple, limited scale |
| **Message Queue** | Reliable delivery |

### Redis Pub/Sub

```typescript
import { createClient } from 'redis';

const pub = createClient();
const sub = createClient();

// Subscribe to channel
sub.subscribe('updates', (message) => {
  broadcastToLocalClients(JSON.parse(message));
});

// Publish from any server
pub.publish('updates', JSON.stringify(data));
```

---

## ✅ Checklist

- [ ] **Reconnection implemented?** (exponential backoff)
- [ ] **Heartbeat configured?** (detect dead connections)
- [ ] **Message format defined?** (consistent protocol)
- [ ] **Error handling robust?** (graceful degradation)
- [ ] **Scaling considered?** (pub/sub if multi-server)
- [ ] **State synchronized?** (initial sync on connect)

---

> **Remember:** Start with SSE. Only use WebSocket if you need bidirectional communication. Only use CRDTs if you need conflict resolution.

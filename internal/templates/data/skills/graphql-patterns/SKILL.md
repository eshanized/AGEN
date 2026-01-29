---
name: graphql-patterns
description: GraphQL schema design, resolvers, federation, and best practices.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# GraphQL Patterns

> Principles for designing and implementing GraphQL APIs.
> **Schema-first thinking.**

---

## 1. Schema Design Principles

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| **Types** | PascalCase | `User`, `OrderItem` |
| **Fields** | camelCase | `firstName`, `createdAt` |
| **Enums** | SCREAMING_SNAKE | `ORDER_STATUS`, `PENDING` |
| **Arguments** | camelCase | `filter`, `orderBy` |

### Type Design

| Principle | Implementation |
|-----------|----------------|
| **Client-centric** | Design for UI needs |
| **Nullable by default** | Explicit non-null (!) |
| **Single responsibility** | One purpose per type |
| **Use interfaces** | Common fields abstracted |

---

## 2. Query Design

### Pagination

```graphql
# Relay-style connections (recommended)
type Query {
  users(first: Int, after: String): UserConnection!
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
}

type UserEdge {
  cursor: String!
  node: User!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Filtering & Sorting

```graphql
input UserFilter {
  status: UserStatus
  createdAfter: DateTime
  search: String
}

input UserOrderBy {
  field: UserOrderField!
  direction: OrderDirection!
}

type Query {
  users(
    filter: UserFilter
    orderBy: UserOrderBy
    first: Int
    after: String
  ): UserConnection!
}
```

---

## 3. Mutation Design

### Input Types

```graphql
# Separate input for each mutation
input CreateUserInput {
  email: String!
  name: String!
}

input UpdateUserInput {
  name: String
  avatar: String
}

type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUser(id: ID!, input: UpdateUserInput!): UpdateUserPayload!
}
```

### Payload Types

```graphql
# Return type with potential errors
type CreateUserPayload {
  user: User
  errors: [UserError!]!
}

type UserError {
  field: String
  message: String!
  code: String!
}
```

---

## 4. N+1 Problem

### The Problem

```
Query: users { orders { items } }
→ 1 query for users
→ N queries for orders (one per user)
→ N*M queries for items (one per order)
```

### Solution: DataLoader

```typescript
// Create per-request DataLoader
const orderLoader = new DataLoader(async (userIds) => {
  const orders = await db.orders.findMany({
    where: { userId: { in: userIds } }
  });
  // Return in same order as userIds
  return userIds.map(id => orders.filter(o => o.userId === id));
});

// Use in resolver
User: {
  orders: (user, _, { loaders }) => loaders.orderLoader.load(user.id)
}
```

---

## 5. Error Handling

### Error Types

| Type | Use |
|------|-----|
| **UserError** | Expected errors (validation, not found) |
| **GraphQL Errors** | Unexpected errors (throw) |

### Pattern

```graphql
type Mutation {
  updateUser(id: ID!, input: UpdateUserInput!): UpdateUserPayload!
}

type UpdateUserPayload {
  user: User
  userErrors: [UserError!]!
}
```

---

## 6. Security

### Authorization

```typescript
// Field-level authorization
const resolvers = {
  User: {
    email: (user, _, { currentUser }) => {
      if (user.id !== currentUser.id) return null;
      return user.email;
    }
  }
};
```

### Query Complexity

```typescript
// Limit query depth and complexity
const server = new ApolloServer({
  validationRules: [
    depthLimit(10),
    queryComplexity({ maximumComplexity: 1000 })
  ]
});
```

---

## 7. Federation (Multi-Service)

### Entity Definition

```graphql
# Users service
type User @key(fields: "id") {
  id: ID!
  name: String!
}

# Orders service
extend type User @key(fields: "id") {
  id: ID! @external
  orders: [Order!]!
}
```

---

## ✅ Checklist

- [ ] **Schema follows conventions?** (naming, types)
- [ ] **Pagination implemented?** (connections)
- [ ] **N+1 solved?** (DataLoader)
- [ ] **Error handling consistent?** (payload patterns)
- [ ] **Authorization in place?** (field-level)
- [ ] **Complexity limited?** (depth, cost)
- [ ] **Documented?** (descriptions on types)

---

> **Remember:** GraphQL is a contract with clients. Design the schema for them, not your database.

# Example queries — allowed vs blocked

Playground: `http://localhost:3003`. Paste the header into the **Headers** tab, the query into
the editor, and click **Run**. The unique `X-API-Key` per feature keeps each example isolated.

> These use sample subgraphs (`employees`, `products`). Adapt the field names to your own
> schema when integrating.

---

## Feature 1 — Response-size limit  (custom module, limit = 400 bytes)
Header: `{ "X-API-Key": "size-check" }`

**Works (small response → HTTP 200):**
```graphql
query Size1 { employee(id: 1) { id } }
query Size2 { employees { id } }
```
**Blocked (big response → RESPONSE_SIZE_LIMIT_EXCEEDED):**
```graphql
query Size3 { employees { id tag expertise notes } }
query Size4 { employees { id tag expertise notes updatedAt } }
```

---

## Feature 2 — Depth limit  (limit = 6)
Header: `{ "X-API-Key": "depth-check" }`

**Works (shallow → HTTP 200):**
```graphql
query Depth1 { employee(id: 1) { id details { forename surname } } }
query Depth2 { employees { role { title } } }
```
**Blocked (too deep → HTTP 400, "query depth 7 exceeds the max query depth allowed (6)"):**
```graphql
query Depth3 { employees { role { employees { role { employees { role { departments } } } } } } }
query Depth4 { employee(id: 1) { role { employees { role { employees { role { title } } } } } } }
```

---

## Feature 3 — Query complexity  (cost limit = 200)
Header: `{ "X-API-Key": "cost-check" }`

**Works (cheap → HTTP 200):**
```graphql
query Cost1 { employees { id } }
query Cost2 { sharedThings(numOfA: 20, numOfB: 1) { a } }
```
**Blocked (too expensive → HTTP 400, "estimated query cost … exceeds the maximum allowed limit 200"):**
```graphql
query Cost3 { productTypes { __typename } }
query Cost4 { sharedThings(numOfA: 5000, numOfB: 1) { a } }
```

Note: `sharedThings` cost equals its `numOfA` argument — an easy dial to cross the limit.

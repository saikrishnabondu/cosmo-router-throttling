# Response-size limit — example queries

Playground: `http://localhost:3003`. Paste the header into the **Headers** tab, the query
into the editor, and click **Run**. Limit = 400 bytes.

Header: `{ "X-API-Key": "size-check" }`

> These use sample subgraphs (`employees`, `products`). Adapt the field names to your own
> schema when integrating.

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

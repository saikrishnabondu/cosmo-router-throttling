# Query cost limit — example queries

Playground: `http://localhost:3003`. Paste the header into the **Headers** tab, the query
into the editor, and click **Run**. Cost limit = 200.

Header: `{ "X-API-Key": "cost-check" }`

> These use sample subgraphs (`employees`, `products`). Adapt the field names to your own
> schema when integrating.

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

# Query depth limit — example queries

Playground: `http://localhost:3003`. Paste the header into the **Headers** tab, the query
into the editor, and click **Run**. Limit = 6 levels.

Header: `{ "X-API-Key": "depth-check" }`

> These use sample subgraphs (`employees`, `products`). Adapt the field names to your own
> schema when integrating.

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

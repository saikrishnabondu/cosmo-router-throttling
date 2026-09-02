# N+1 Hydration — example query

Router playground: `http://localhost:3003`. Paste the query in the editor, click **Run**.
To *see* the batched execution, set the **Headers** tab to:
```json
{ "X-WG-TRACE": "true" }
```
and expand `extensions` in the response — you'll see each subgraph fetch and its timing.

A list of employees, where each employee's `products` is hydrated from the **products**
service and joined back by `id`.

```graphql
query N1Hydration { employees { id tag products } }
```

**What you get** (shape):
```json
{ "data": { "employees": [
  { "id": 1, "tag": "...", "products": ["COSMO","SDK"] },
  { "id": 2, "tag": "...", "products": ["CONSULTANCY"] }
] } }
```

**What the router did:**
- 1 call to the **employees** service → the list (`id`, `tag`, and each employee's key).
- **1 batched call** to the **products** service with *all* the ids at once → every employee's
  `products`.
- Joined each `products` back to its employee by `id`.

So it's **1 + 1 calls**, not 1 + N. In the `X-WG-TRACE` output you'll see a single
`_entities` fetch to the products service carrying many representations.

> Contrast: `query { employees { id tag } }` touches only the employees service — no hydration.
> Adding the `products` field is what triggers the batched hydration from the other service.

## Proving the "N+1 avoided" point
Run `N1Hydration` for the full employee list and check the trace — there is **one** fetch to
the products service regardless of how many employees came back. That single batched call is
the whole point of hydration.

# Scatter-Gather Composition — example query

Router playground: `http://localhost:3003`. Paste the query in the editor, click **Run**.
To *see* the parallel execution, set the **Headers** tab to:
```json
{ "X-WG-TRACE": "true" }
```
and expand `extensions` in the response — you'll see each subgraph fetch and its timing.

One query that pulls from **two independent services at the same time** and returns a single
combined response.

```graphql
query ScatterGather {
  employees { id }                          # employees service
  sharedThings(numOfA: 3, numOfB: 1) { a }  # products service
}
```

**What you get** (shape):
```json
{ "data": {
  "employees": [ { "id": 1 }, { "id": 2 } ],
  "sharedThings": [ { "a": "a-0" }, { "a": "a-1" }, { "a": "a-2" } ]
} }
```

**What the router did:**
- Saw two independent root fields owned by different services.
- Fetched **both at the same time** (parallel group).
- **Stitched** them into one response.

Total time ≈ the slower of the two fetches, not the sum. In the `X-WG-TRACE` output the two
fetches **overlap** in time.

> In a real dashboard this fans out to many services (accounts, cards, fraud, rewards…). This
> setup has two services running, so it shows a 2-way fan-out; the pattern is identical at
> scale.

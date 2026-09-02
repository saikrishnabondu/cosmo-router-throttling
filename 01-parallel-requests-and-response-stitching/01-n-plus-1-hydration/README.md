# N+1 Hydration

After fetching a **list** of objects, the details for each object are fetched in **one batched,
parallel call** to another service and **joined back by key (`id`)** — avoiding one call per
item (the "N+1" problem).

No custom code and no config switch. This is what the Cosmo Router's **query planner** does by
default, as long as the schema is federated correctly (entities declared with `@key`).

## Files
| File | What it is |
| --- | --- |
| `queries.md` | The query that triggers it, expected shape, and how to see the batched call |

## How it works

**The naive (bad) way — the N+1 problem:** fetch the list (1 call), then make **one call per
item** to get each detail (N calls). For 50 employees that's 51 calls — slow, and it hammers
the backend.

**How Cosmo does it instead:**
1. **Fetch the list** from the owning subgraph — e.g. `employees { id }` from the *employees*
   service. Each row carries its `@key` (here `id`).
2. **Collect the keys** of every object into one list of "entity representations".
3. **One batched call** to the other service using the federation `_entities` query, passing
   **all keys at once** — e.g. hydrate `products` for all 50 employees in a single request to
   the *products* service.
4. **Join back by key** — the stitcher attaches each returned `products` to the employee with
   the matching `id`.

So N+1 becomes **1 + 1**: one call for the list, one batched call for all the details.

The key enabler is the shared entity:
```graphql
# employees service owns the entity
type Employee @key(fields: "id") { id: Int!  tag: String! }
# products service extends the SAME entity, adding a field
type Employee @key(fields: "id") { id: Int!  products: [ProductName!]! }
```
Because both declare `@key(fields: "id")`, the router knows how to fetch and join them. Full
schema excerpt: [`../schema/federation-keys.graphqls`](../schema/federation-keys.graphqls).

Run steps: [`../docs/setup.md`](../docs/setup.md).

# Parallel Requests and Response Stitching

**Status: partially implemented.** N+1 Hydration and Scatter-Gather Composition are in; more may follow.

Two federation patterns, each in its own folder:

| # | Feature | What it does |
|---|---|---|
| 1 | [N+1 Hydration](01-n-plus-1-hydration/) | After fetching a **list** of objects, the details for each object are fetched in **one batched, parallel call** to another service and **joined back by key (`id`)** — avoiding one call per item (the "N+1" problem). |
| 2 | [Scatter-Gather Composition](02-scatter-gather-composition/) | **Multiple independent services** are called **at the same time** and their responses are **combined into one response** — ideal for dashboards that depend on many backends. |

## These are *automatic* router behaviours

Unlike the throttling features, you do **not** write code or flip a config switch to "enable"
these. They are what the Cosmo Router's **query planner** does by default, **as long as the
schema is federated correctly** (entities declared with `@key`). So the work here is:

1. A **federated schema** where entities are shared across subgraphs via `@key`
   (each feature folder ships its own `schema/federation-keys.graphqls`).
2. **Queries** that trigger each pattern (in each feature folder).
3. Proof that the router executes them in **parallel** and **stitches** the result.

The `Employee` entity is shared between the **employees** and **products** subgraphs (employees
owns `id`; products adds `Employee.products`), which is a real N+1 hydration relationship.

## Where each pattern sits in the request lifecycle
1. Request arrives → validate → plan.
2. Planner emits: **parallel groups** (scatter-gather), **sequence nodes** (dependencies), and
   **batched entity fetches** (hydration).
3. Fetches execute (parallel where independent).
4. **Response stitcher** merges: root merge, entity merge by `@key`, nested placement,
   list-slot merge, plus partial-failure handling.
5. One response returned.

## Proving the parallelism
Send any query with the header `{ "X-WG-TRACE": "true" }`. The response's `extensions` include
tracing that shows each subgraph fetch and its timing — the hydration fetch is a **single
batched call**, and the scatter-gather fetches **overlap in time** rather than running one
after another. (The same is visible as OpenTelemetry spans if a collector is attached.)

## Layout

Each feature folder is **self-contained** — its own schema excerpt, router config, composed
supergraph, start scripts, query, and run steps:

```
01-n-plus-1-hydration/            the batched entity fetch
02-scatter-gather-composition/    the parallel root fetch

  ...each containing:
    README.md                       how the pattern works
    queries.md                      the query that triggers it + expected result
    setup.md                        how to run
    schema/federation-keys.graphqls the @key relationships (reference)
    router-config/                  clean router config, composition manifest, composed supergraph
    scripts/                        start infra, subgraphs, and the router (Windows PowerShell)
```

Both run on a **clean router config** with **no throttling features**, using the stock
`router.exe`, and both listen on `http://localhost:3003` — so run one at a time.

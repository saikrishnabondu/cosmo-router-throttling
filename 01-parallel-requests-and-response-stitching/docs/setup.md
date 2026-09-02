# How to run (standalone — no throttling involved)

This use case is **separate from Response Throttling**. It runs on its **own clean router
config** (`router-config/router.config.yaml`) with **no throttling features**, so nothing
interferes with the parallel-fetch / stitching demo. It uses the stock `router.exe`.

The only thing shared with other demos is the **backend services** (the `employees` and
`products` subgraphs from `C:\cosmo\demo`) — those are just data services, not throttling.

## Prerequisites
- Docker infra up (Redis + NATS) so the schema builds:
  ```
  cd C:\cosmo ; docker compose up -d redis nats
  ```
- Go on PATH, and the `C:\cosmo` repo present (for the subgraph services and router.exe).

## Start it (from this folder's `scripts`)
```
powershell -File "scripts/01-start-subgraphs.ps1"
powershell -File "scripts/02-start-router.ps1"
```
- `01` starts employees (4101) + products (4004). If they're already running you'll see a
  harmless bind error — ignore it.
- `02` frees port 3003 and starts the **clean** router (no throttling) against this config.

## Then
1. Open `http://localhost:3003`.
2. Run the queries in [`01-n-plus-1-hydration/queries.md`](../01-n-plus-1-hydration/queries.md) and [`02-scatter-gather-composition/queries.md`](../02-scatter-gather-composition/queries.md).
3. Add the header `{ "X-WG-TRACE": "true" }` to see the batched hydration call and the
   overlapping parallel fetches in the response `extensions`.

Because this router has no response-size / cost / depth limits, the hydration query returns the
full result — the throttling features from the other project can't get in the way here.

## What "implementation" means for this use case
No custom code, no throttling config. The patterns are the router's **default federation
behaviour**, produced by the query planner from the **federated schema** (the shared `@key` on
`Employee`). The deliverable is: the clean router + composed graph here, the federated schema
that exhibits the patterns (`schema/federation-keys.graphqls`), the queries that trigger them
(one `queries.md` per feature folder), and the explanation of how the planner executes them
(each feature's README).

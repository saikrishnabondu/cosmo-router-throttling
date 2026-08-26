# Response Throttling

**Status: implemented.**

Bounding response size, query depth, and query cost at the router, in front of the backend
services — so the services themselves need no changes.

| # | Feature | Type | What it does |
|---|---------|------|--------------|
| 1 | [Response-size limit](01-response-size-limit/) | **Custom Go module** (Cosmo has no native option) | Rejects responses whose body is larger than a configured byte limit |
| 2 | [Query depth limit](02-query-depth-limit/) | Native Cosmo config | Rejects queries nested deeper than a set number of levels |
| 3 | [Query cost limit](03-query-cost-limit/) | Native Cosmo config | Rejects queries whose estimated cost exceeds a set budget |

More throttling features may be added here over time.

Each feature folder is self-contained: its own `router.config.yaml`, its own `queries.md`, and
a README explaining the internals. Build and run steps shared by all three live in
[`docs/setup.md`](docs/setup.md).

## Where each feature runs in the request lifecycle
1. Request arrives at the router.
2. Middleware runs.
3. Parse → validate → **plan** → **depth check + cost analysis** happen here (native Cosmo).
4. Execute → call subgraph → **`OnOriginResponse`** runs → **response-size module** here.
5. Router merges results → returns to the client.

## Running more than one at a time

The three configs are separate so each feature can be run and demonstrated on its own. To
enforce several at once, merge the feature-specific blocks into a single `router.config.yaml` —
they occupy different keys (`modules`, `security.complexity_limits`, `security.cost_control`)
and do not conflict.

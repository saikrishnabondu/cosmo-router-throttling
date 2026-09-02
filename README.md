# Federated API Demo

Federated API use cases on a [WunderGraph Cosmo](https://cosmo-docs.wundergraph.com/)
federated GraphQL graph. Every use case is handled in the **router**, in front of the backend
services, so the services themselves need no changes.

Each use case lives in its own folder.

| # | Use case | What it covers | Status |
|---|---|---|---|
| 1 | [Parallel Requests and Response Stitching](01-parallel-requests-and-response-stitching/) | Fanning one client query out to multiple subgraphs concurrently and assembling a single response | [**Partially implemented**](#status-notes) |
| 2 | [API Chaining](02-api-chaining/) | Feeding the output of one subgraph call into the next as an ordered dependency | Not planned |
| 3 | [Protocol Conversion](03-protocol-conversion/) | Fronting non-GraphQL backends (REST, gRPC) behind one federated GraphQL surface | Not planned |
| 4 | [Schema Evolution](04-schema-evolution/) | Changing the graph over time — deprecation, additive change, versioning — without breaking clients | Not planned |
| 5 | [API Provisioning](05-api-provisioning/) | Registering subgraphs, composing the schema, and rolling out router configuration | Not planned |
| 6 | [Response Throttling](06-response-throttling/) | Bounding response size, query depth, and query cost at the router | [**Partially implemented**](#status-notes) |

Statuses are summarised below.

## Status notes

Click a section to expand.

<details>
<summary><b>1 — Parallel Requests and Response Stitching — partially implemented</b></summary>

Implemented:

1. [N+1 Hydration](01-parallel-requests-and-response-stitching/01-n-plus-1-hydration/)
2. [Scatter-Gather Composition](01-parallel-requests-and-response-stitching/02-scatter-gather-composition/)

</details>

<details>
<summary><b>6 — Response Throttling — partially implemented</b></summary>

Implemented:

1. [Response-size limit](06-response-throttling/01-response-size-limit/)
2. [Query depth limit](06-response-throttling/02-query-depth-limit/)
3. [Query cost limit](06-response-throttling/03-query-cost-limit/)

</details>

<details>
<summary><b>2–5 — Not planned</b></summary>

These folders exist as placeholders so the structure is visible. There is no scheduled work on
them; each holds a README stub naming the use case. If one is picked up later, its status here
moves to planned or implemented.

</details>

---

## Contributors

See [`CONTRIBUTORS.md`](CONTRIBUTORS.md).

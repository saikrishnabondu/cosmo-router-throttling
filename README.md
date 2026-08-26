# Federated API Demo

Federated API use cases on a [WunderGraph Cosmo](https://cosmo-docs.wundergraph.com/)
federated GraphQL graph. Every use case is handled in the **router**, in front of the backend
services, so the services themselves need no changes.

Each use case lives in its own folder.

| # | Use case | What it covers | Status |
|---|---|---|---|
| 1 | [Parallel Requests and Response Stitching](01-parallel-requests-and-response-stitching/) | Fanning one client query out to multiple subgraphs concurrently and assembling a single response | Not planned |
| 2 | [API Chaining](02-api-chaining/) | Feeding the output of one subgraph call into the next as an ordered dependency | Not planned |
| 3 | [Protocol Conversion](03-protocol-conversion/) | Fronting non-GraphQL backends (REST, gRPC) behind one federated GraphQL surface | Not planned |
| 4 | [Schema Evolution](04-schema-evolution/) | Changing the graph over time — deprecation, additive change, versioning — without breaking clients | Not planned |
| 5 | [API Provisioning](05-api-provisioning/) | Registering subgraphs, composing the schema, and rolling out router configuration | Not planned |
| 6 | [Response Throttling](06-response-throttling/) | Bounding response size, query depth, and query cost at the router | [**Partially implemented**](#status-notes) |

Statuses are summarised below.

## Status notes

Click a section to expand.

<details>
<summary><b>6 — Response Throttling — partially implemented</b></summary>

Three features are built, each in its own self-contained folder with config, example queries,
and an explanation of the internals:

| Feature | Type | State |
| --- | --- | --- |
| [Response-size limit](06-response-throttling/01-response-size-limit/) | Custom Go module (Cosmo has no native option) | Working, with one caveat below |
| [Query depth limit](06-response-throttling/02-query-depth-limit/) | Native Cosmo config | Working |
| [Query cost limit](06-response-throttling/03-query-cost-limit/) | Native Cosmo config | Working |

**Caveat on the response-size limit.** Because it enforces at the subgraph-response stage
(`OnOriginResponse`), the client sees the rejection as a GraphQL error inside an HTTP 200
envelope carrying the `RESPONSE_SIZE_LIMIT_EXCEEDED` code — not a bare HTTP 413. The size cap
itself works as intended.

**Why "partially".** Response throttling covers more ground than these three limits; further
features may be added to this folder over time.

Build and run steps: [`06-response-throttling/docs/setup.md`](06-response-throttling/docs/setup.md).

</details>

<details>
<summary><b>1–5 — Not planned</b></summary>

These folders exist as placeholders so the structure is visible. There is no scheduled work on
them; each holds a README stub naming the use case. If one is picked up later, its status here
moves to planned or implemented.

</details>

---

## Contributors

See [`CONTRIBUTORS.md`](CONTRIBUTORS.md).

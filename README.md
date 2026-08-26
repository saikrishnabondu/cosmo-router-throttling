# Federated API Use Cases — WunderGraph Cosmo

Federated API use cases on a [WunderGraph Cosmo](https://cosmo-docs.wundergraph.com/)
federated GraphQL graph. Every use case is handled in the **router**, in front of the backend
services, so the services themselves need no changes.

Each use case lives in its own folder.

| # | Use case | What it covers | Status |
|---|---|---|---|
| 1 | [Parallel Requests and Response Stitching](01-parallel-requests-and-response-stitching/) | Fanning one client query out to multiple subgraphs concurrently and assembling a single response | Planned |
| 2 | [API Chaining](02-api-chaining/) | Feeding the output of one subgraph call into the next as an ordered dependency | Planned |
| 3 | [Protocol Conversion](03-protocol-conversion/) | Fronting non-GraphQL backends (REST, gRPC) behind one federated GraphQL surface | Planned |
| 4 | [Schema Evolution](04-schema-evolution/) | Changing the graph over time — deprecation, additive change, versioning — without breaking clients | Planned |
| 5 | [API Provisioning](05-api-provisioning/) | Registering subgraphs, composing the schema, and rolling out router configuration | Planned |
| 6 | [Response Throttling](06-response-throttling/) | Bounding response size, query depth, and query cost at the router | **Implemented** |

More use cases may be added over time.

## Implemented today

[**Response Throttling**](06-response-throttling/) covers three features, each in its own
self-contained folder with config, example queries, and an explanation of the internals:

1. [Response-size limit](06-response-throttling/01-response-size-limit/) — custom Go module
   (Cosmo has no native option)
2. [Query depth limit](06-response-throttling/02-query-depth-limit/) — native Cosmo config
3. [Query cost limit](06-response-throttling/03-query-cost-limit/) — native Cosmo config

Build and run steps: [`06-response-throttling/docs/setup.md`](06-response-throttling/docs/setup.md).

---

## Contributors

See [`CONTRIBUTORS.md`](CONTRIBUTORS.md).

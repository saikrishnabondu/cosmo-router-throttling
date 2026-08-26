# Cosmo Router Throttling — 3 Features

Response throttling for a [WunderGraph Cosmo](https://cosmo-docs.wundergraph.com/) federated
GraphQL graph, focused on **three features**:

| # | Feature | Type | What it does |
|---|---------|------|--------------|
| 1 | **Response-size limit** | **Custom Go module** (Cosmo has no native option) | Rejects responses whose body is larger than a configured byte limit |
| 2 | **Depth limit** | Cosmo native (config) | Rejects queries nested deeper than a configured number of levels |
| 3 | **Query complexity** | Cosmo native (config) | Prices each query from `@cost` / `@listSize` directives and rejects expensive ones before execution |

Every rule is enforced in the **router**, in front of the backend services, so the services
themselves need no changes.

---

## Repository layout

```
router-config/
  router.config.yaml     # depth limit + cost control + the responseSizeLimit module
custom-module/
  main.go                # custom router entrypoint (imports the module)
  modules/
    response_size.go     # Feature 1: the response-size limit module (Go)
examples/
  queries.md             # copy-paste "work" vs "blocked" queries for all 3 features
docs/
  setup.md               # how to build and run
  how-it-works.md        # internals: how cost is calculated & how the module is triggered
```

---

## Quick start

Response-size (Feature 1) needs a **custom router build** because it's a Go module; depth and
cost (Features 2 & 3) work on any Cosmo router. See [`docs/setup.md`](docs/setup.md) for the
full steps. In short:

```bash
# from a checkout of the cosmo router repo, with these module files placed under cmd/throttle-router
go build -o throttle-router.exe ./cmd/throttle-router
./throttle-router.exe -config router.config.yaml
```

Then open the playground at `http://localhost:3003` and run the queries in
[`examples/queries.md`](examples/queries.md).

---

## How each feature works (summary)

### 1. Response-size limit (custom module)
Cosmo only limits the **request** direction (`max_request_body_size`); there is no
response-size option in its config. This module implements the `OnOriginResponse` hook — it
runs **after a subgraph responds**, measures the response body, and if it exceeds
`max_response_bytes` it returns a GraphQL error with code `RESPONSE_SIZE_LIMIT_EXCEEDED`.
Otherwise it passes the response through unchanged. Code: [`custom-module/modules/response_size.go`](custom-module/modules/response_size.go).

```yaml
modules:
  responseSizeLimit:
    max_response_bytes: 400
```

### 2. Depth limit (native config)
Rejects queries nested deeper than the limit, at parse time — blocks recursive-query attacks.

```yaml
security:
  complexity_limits:
    mode: enforce
    depth: { enabled: true, limit: 6 }
```

### 3. Query complexity (native config)
The router computes an estimated cost from `@cost` weights and `@listSize` hints in the schema
and rejects anything over `max_estimated_limit`, before any backend call.

```yaml
security:
  cost_control:
    enabled: true
    mode: enforce
    max_estimated_limit: 200
```

See [`docs/how-it-works.md`](docs/how-it-works.md) for the exact cost math and how the router
dispatches to the custom module.

---

## Status

Verified against a local Cosmo Router with sample subgraphs. Each feature is proven
with a query that succeeds (HTTP 200) and a query that is blocked (HTTP 400 / error code).

---

## Contributors

See [`CONTRIBUTORS.md`](CONTRIBUTORS.md).

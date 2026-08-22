# How it works — internals

## Request lifecycle (where each feature runs)
1. Request arrives at the router.
2. Middleware runs.
3. Parse → validate → **plan** → **depth check + cost analysis** happen here (native Cosmo).
4. Execute → call subgraph → **`OnOriginResponse`** runs → **response-size module** here.
5. Router merges results → returns to the client.

---

## Query complexity — how the cost is calculated

Native Cosmo feature (enabled via `security.cost_control`). Runs during planning, **before
execution**. Rules:

- Scalar/enum leaf fields → cost **0**.
- Object/interface/union fields → base **1**.
- `@cost(weight: N)` → sets that field/type cost to **N**.
- `@listSize(assumedSize: N)` → multiply the field's contents by **N**.
- `@listSize(slicingArguments: ["numOfA"])` → list size comes from the client's argument.
- Union → uses the most expensive member.

Worked examples (limit 200), using the sample schema:

| Query | Calculation | Cost | Result |
|---|---|---|---|
| `productTypes { __typename }` | 50 (assumedSize) × 8 (`Cosmo @cost weight`) | 400 | rejected |
| `sharedThings(numOfA: 20)` | 20 × 1 | 20 | allowed |
| `sharedThings(numOfA: 5000)` | 5000 × 1 | 5000 | rejected |
| `employees { id }` | 50 × 1 | 50 | allowed |
| `employee(id: 1) { id }` | field weight 5 + arg weight 2 | 7 | allowed |

Rejection message: `"The estimated query cost 400 exceeds the maximum allowed limit 200"` (HTTP 400).

---

## Depth limit

Native Cosmo feature (`security.complexity_limits.depth`). At parse time the router counts the
nesting depth of the query. Each `{ … }` is one level. If the depth exceeds the limit (6), the
query is rejected with `"The query depth 7 exceeds the max query depth allowed (6)"` (HTTP 400) —
before execution.

---

## Response-size limit — how the custom module is triggered

**Not matched by file name — matched by Go interface.** Three steps:

1. **Register** at startup:
   ```go
   func init() { core.RegisterModule(&ResponseSizeLimitModule{}) }
   ```
2. **Declare the hook** by implementing an interface. The module has an `OnOriginResponse`
   method, which is the single method of `core.EnginePostOriginHandler`.
3. **The router wires it up.** At startup it checks every module with a type assertion — this is
   the actual dispatch in `router/core/router.go`:
   ```go
   if handler, ok := moduleInstance.(EnginePostOriginHandler); ok {
       r.postOriginHandlers = append(r.postOriginHandlers, handler.OnOriginResponse)
   }
   ```
   Plain English: *"Does this module implement the after-response hook? If yes, add its
   `OnOriginResponse` to the list that runs after every subgraph response."*

So the module runs at exactly that point because that's the interface it implements. Settings
reach it via the module **ID** (`responseSizeLimit`), which the router matches to the
`modules.responseSizeLimit` section in the config.

### What OnOriginResponse does
```go
body, _ := io.ReadAll(resp.Body)          // read the response the backend produced
if len(body) > m.MaxResponseBytes {        // too big?
    return <error with RESPONSE_SIZE_LIMIT_EXCEEDED>   // reject
}
resp.Body = io.NopCloser(bytes.NewReader(body))
return resp                                // within limit -> pass through unchanged
```

Honest detail: because it enforces at the subgraph-response stage, the client sees the
rejection as a GraphQL error inside an HTTP 200 envelope (carrying the
`RESPONSE_SIZE_LIMIT_EXCEEDED` code), not a bare HTTP 413. The size cap itself works as intended.

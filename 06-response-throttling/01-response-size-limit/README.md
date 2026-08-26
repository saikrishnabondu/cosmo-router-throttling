# Response-size limit

**Status: implemented.** Custom Go module — Cosmo has no native option for this.

Rejects responses whose body is larger than a configured byte limit
(`max_response_bytes`, default 400 here).

## Files
| File | What it is |
| --- | --- |
| `main.go` | Custom router entrypoint that imports the module |
| `modules/response_size.go` | The module itself |
| `router.config.yaml` | Standalone config for this feature |
| `queries.md` | Allowed vs blocked example queries |

## How it is triggered

**Not matched by file name — matched by Go interface.** Three steps:

1. **Register** at startup:
   ```go
   func init() { core.RegisterModule(&ResponseSizeLimitModule{}) }
   ```
2. **Declare the hook** by implementing an interface. The module has an `OnOriginResponse`
   method, which is the single method of `core.EnginePostOriginHandler`.
3. **The router wires it up.** At startup it checks every module with a type assertion — this
   is the actual dispatch in `router/core/router.go`:
   ```go
   if handler, ok := moduleInstance.(EnginePostOriginHandler); ok {
       r.postOriginHandlers = append(r.postOriginHandlers, handler.OnOriginResponse)
   }
   ```
   Plain English: *"Does this module implement the after-response hook? If yes, add its
   `OnOriginResponse` to the list that runs after every subgraph response."*

Settings reach it via the module **ID** (`responseSizeLimit`), which the router matches to the
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
`RESPONSE_SIZE_LIMIT_EXCEEDED` code), not a bare HTTP 413. The size cap itself works as
intended.

Build and run steps: [`../docs/setup.md`](../docs/setup.md).

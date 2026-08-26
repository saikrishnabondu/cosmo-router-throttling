# Query depth limit

**Status: implemented.** Native Cosmo feature — works on any stock router build, no custom
compile needed.

Rejects queries nested deeper than the configured limit (6 here), via
`security.complexity_limits.depth`.

## Files
| File | What it is |
| --- | --- |
| `router.config.yaml` | Standalone config for this feature |
| `queries.md` | Allowed vs blocked example queries |

## How it works

At parse time the router counts the nesting depth of the query. Each `{ … }` is one level. If
the depth exceeds the limit, the query is rejected with
`"The query depth 7 exceeds the max query depth allowed (6)"` (HTTP 400) — **before execution**,
so no subgraph is ever called.

The same config block also caps `total_fields` (60) and `root_fields` (5).

Run steps: [`../docs/setup.md`](../docs/setup.md).

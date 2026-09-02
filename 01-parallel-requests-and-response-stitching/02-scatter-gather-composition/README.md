# Scatter-Gather Composition

**Status: implemented.**

**Multiple independent services** are called **at the same time** and their responses are
**combined into one response** — ideal for dashboards that depend on many backends.

No custom code and no config switch. This is what the Cosmo Router's **query planner** does by
default when a query touches independent root fields owned by different subgraphs.

## Files
| File | What it is |
| --- | --- |
| `queries.md` | The query that triggers it, expected shape, and how to see the parallelism |

## How it works

1. The planner sees the query touches multiple root fields owned by **different** subgraphs
   with **no dependency** between them.
2. It puts those fetches in a **parallel group** and dispatches them **at the same time**
   (concurrent goroutines), bounded by `max_concurrent_resolvers` and per-fetch timeouts.
3. The **response stitcher** merges the separate results into one JSON (a "root merge"),
   applying field security and partial-failure handling.

So instead of the client calling service A, then B, then C (a waterfall), the router calls them
**together** and returns one combined response. Total time ≈ the **slowest** service, not the
**sum**.

Schema excerpt showing the independent roots:
[`schema/federation-keys.graphqls`](schema/federation-keys.graphqls).

Run steps: [`setup.md`](setup.md).

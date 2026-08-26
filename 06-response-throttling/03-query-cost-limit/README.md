# Query cost limit (complexity)

**Status: implemented.** Native Cosmo feature — works on any stock router build, no custom
compile needed.

Rejects queries whose *estimated* cost exceeds the limit (200 here), via
`security.cost_control`. Runs during planning, **before execution**.

## Files
| File | What it is |
| --- | --- |
| `router.config.yaml` | Standalone config for this feature |
| `queries.md` | Allowed vs blocked example queries |

## How the cost is calculated

- Scalar/enum leaf fields → cost **0**.
- Object/interface/union fields → base **1**.
- `@cost(weight: N)` → sets that field/type cost to **N**.
- `@listSize(assumedSize: N)` → multiply the field's contents by **N**.
- `@listSize(slicingArguments: ["numOfA"])` → list size comes from the client's argument.
- Union → uses the most expensive member.

Worked examples (limit 200), using the sample schema:

| Query | Calculation | Cost | Result |
|---|---|---|---|
| `productTypes { __typename }` | 50 (assumedSize) × 8 (`@cost` weight) | 400 | rejected |
| `sharedThings(numOfA: 20)` | 20 × 1 | 20 | allowed |
| `sharedThings(numOfA: 5000)` | 5000 × 1 | 5000 | rejected |
| `employees { id }` | 50 × 1 | 50 | allowed |
| `employee(id: 1) { id }` | field weight 5 + arg weight 2 | 7 | allowed |

Rejection message: `"The estimated query cost 400 exceeds the maximum allowed limit 200"`
(HTTP 400).

For this to reject anything, the schema must carry `@cost` / `@listSize` directives.

Run steps: [`../docs/setup.md`](../docs/setup.md).

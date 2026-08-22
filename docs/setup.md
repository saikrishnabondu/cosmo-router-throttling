# Setup — build and run

## Prerequisites
- Go (same version used to build the Cosmo Router)
- The [wgc](https://www.npmjs.com/package/wgc) CLI (`npm install -g wgc`) to compose the schema
- A checkout of the WunderGraph Cosmo repo (for the router source and sample subgraphs)
- Redis + NATS if you also enable rate limiting (not required for these three features)

## 1. Place the custom module in the router repo
Copy this repo's `custom-module/` into the Cosmo router source as a command:

```
cosmo/router/cmd/throttle-router/main.go              <- custom-module/main.go
cosmo/router/cmd/throttle-router/modules/response_size.go  <- custom-module/modules/response_size.go
```

The import path in `main.go` (`.../cmd/throttle-router/modules`) must match that location.

## 2. Build the custom router
```
cd cosmo/router
go build -o throttle-router.exe ./cmd/throttle-router
```

## 3. Compose the schema
From wherever `graph.yaml` for your subgraphs lives:
```
wgc router compose -i graph.yaml -o config.json
```
Put `config.json` next to `router.config.yaml` (or update `execution_config.file.path`).

## 4. Run
```
./throttle-router.exe -config router.config.yaml
```
On startup the log should show `responseSizeLimit module active`. The playground opens at
`http://localhost:3003`.

## Notes
- **Depth limit** and **query complexity** work on *any* Cosmo router — no custom build needed.
- **Response-size limit** only works on the custom build, because it's a compiled Go module.
- `router.config.yaml` is read only at **startup** — after changing a value (e.g.
  `max_response_bytes`), restart the router for it to take effect.
- For query-complexity to reject expensive queries, the schema must carry `@cost` / `@listSize`
  directives.

package main

import (
	routercmd "github.com/wundergraph/cosmo/router/cmd"

	// Import the custom module so its init() registers it with the router.
	_ "github.com/wundergraph/cosmo/router/cmd/throttle-router/modules"
)

// Custom Cosmo Router entrypoint that adds the response-size limit — a feature
// the stock router has no configuration for.
//
// Build:  go build -o throttle-router.exe ./cmd/throttle-router
// Run:    ./throttle-router.exe -config router.config.yaml
func main() {
	routercmd.Main()
}

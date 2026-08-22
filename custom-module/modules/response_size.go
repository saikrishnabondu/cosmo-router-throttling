package modules

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/wundergraph/cosmo/router/core"
	"go.uber.org/zap"
)

func init() {
	core.RegisterModule(&ResponseSizeLimitModule{})
}

const responseSizeLimitModuleID = "responseSizeLimit"

// ResponseSizeLimitModule caps the size of the response body.
//
// The stock Cosmo Router only limits the REQUEST direction
// (traffic_shaping.router.max_request_body_size / max_header_bytes); there is
// no configuration option anywhere in the Go config structs for a RESPONSE
// size cap. This custom module closes that gap.
//
// It hooks OnOriginResponse (called after a subgraph returns): it reads the
// subgraph response body, and if it is larger than the configured limit it
// replaces it with a GraphQL error. When the body is within the limit it is
// passed through unchanged, so unrelated queries are never affected.
//
// Config (in router.config.yaml):
//
//	modules:
//	  responseSizeLimit:
//	    max_response_bytes: 400
type ResponseSizeLimitModule struct {
	// MaxResponseBytes is the maximum allowed response body size in bytes.
	MaxResponseBytes int `mapstructure:"max_response_bytes"`

	Logger *zap.Logger
}

func (m *ResponseSizeLimitModule) Provision(ctx *core.ModuleContext) error {
	if m.MaxResponseBytes <= 0 {
		return fmt.Errorf("responseSizeLimit: max_response_bytes must be greater than 0")
	}
	m.Logger = ctx.Logger
	m.Logger.Info("responseSizeLimit module active", zap.Int("max_response_bytes", m.MaxResponseBytes))
	return nil
}

func (m *ResponseSizeLimitModule) Cleanup() error { return nil }

// OnOriginResponse inspects the size of each subgraph response.
func (m *ResponseSizeLimitModule) OnOriginResponse(resp *http.Response, ctx core.RequestContext) *http.Response {
	if resp == nil || resp.Body == nil {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader(body))
		return resp
	}

	if len(body) > m.MaxResponseBytes {
		m.Logger.Warn("response exceeded size limit",
			zap.Int("response_bytes", len(body)),
			zap.Int("max_response_bytes", m.MaxResponseBytes),
		)
		errBody := fmt.Sprintf(
			`{"errors":[{"message":"The response size %d bytes exceeds the maximum allowed limit %d bytes","extensions":{"code":"RESPONSE_SIZE_LIMIT_EXCEEDED"}}],"data":null}`,
			len(body), m.MaxResponseBytes,
		)
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		return &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         resp.Proto,
			ProtoMajor:    resp.ProtoMajor,
			ProtoMinor:    resp.ProtoMinor,
			Header:        header,
			Body:          io.NopCloser(bytes.NewReader([]byte(errBody))),
			ContentLength: int64(len(errBody)),
			Request:       resp.Request,
		}
	}

	// Within the limit: restore the body we consumed and pass through unchanged.
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp
}

func (m *ResponseSizeLimitModule) Module() core.ModuleInfo {
	return core.ModuleInfo{
		ID:       responseSizeLimitModuleID,
		Priority: 1,
		New:      func() core.Module { return &ResponseSizeLimitModule{} },
	}
}

// Interface guards
var (
	_ core.EnginePostOriginHandler = (*ResponseSizeLimitModule)(nil)
	_ core.Provisioner             = (*ResponseSizeLimitModule)(nil)
	_ core.Cleaner                 = (*ResponseSizeLimitModule)(nil)
)

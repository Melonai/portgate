package handlers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"portgate"
)

// RequestHandler keeps data relevant to the request handlers.
type RequestHandler struct {
	// Pointer to the global Portgate config, the values of which can change at runtime.
	config *portgate.Config
	// HTTP Client for requesting resources from the destination host.
	client fasthttp.Client
}

func NewRequestHandler(config *portgate.Config) RequestHandler {
	return RequestHandler{
		config: config,
		client: fasthttp.Client{},
	}
}

// HandleRequest handles all types of requests and delegates to more specific handlers.
func (h *RequestHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	path := portgate.ParsePath(string(ctx.Path()))

	if path.DestinationIdentifier == -1 {
		// We were not given a port.

		if path.ResourcePath == "/_portgate" {
			h.handlePortgateRequest(ctx)
		} else {
			// Try to grab actual destination from Referer header.
			refererPath, err := portgate.ParsePathFromReferer(path, string(ctx.Request.Header.Referer()))
			if err != nil || refererPath.DestinationIdentifier == -1 {
				// The referer path also has no destination
				h.handleUnknownRequest(ctx)
			} else {
				// We found the destination from the referer path, so we
				// redirect the user to the Portgate URL they should've requested.

				portgateUrl := fmt.Sprintf("/%d%s", refererPath.DestinationIdentifier, refererPath.ResourcePath)
				ctx.Redirect(portgateUrl, http.StatusTemporaryRedirect)
			}
		}
	} else {
		// We were given a port, so we have to pass the request through to the destination host.

		h.handlePassthroughRequest(ctx, path)
	}
}

// handleUnknownRequest handles any request which could not be processed due to missing
// information.
func (h *RequestHandler) handleUnknownRequest(ctx *fasthttp.RequestCtx) {
	// TODO: Show error page
	ctx.SetStatusCode(http.StatusNotFound)
	_, _ = ctx.WriteString("Unknown request.")
}

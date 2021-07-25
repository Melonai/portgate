package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

// RequestHandler keeps data relevant to the request handlers.
type RequestHandler struct {
	// Pointer to the global Portgate config, the values of which can change at runtime.
	config *Config
	// HTTP Client for requesting resources from the destination host.
	client fasthttp.Client
}

// handleRequest handles all types of requests and delegates to more specific handlers.
func (h *RequestHandler) handleRequest(ctx *fasthttp.RequestCtx) {
	path := ParsePath(string(ctx.Path()))

	if path.DestinationIdentifier == -1 {
		// We were not given a port.

		if path.ResourcePath == "/_portgate" {
			h.handlePortgateRequest(ctx)
		} else {
			// Try to grab actual destination from Referer header.
			refererPath, err := ParsePathFromReferer(path, string(ctx.Request.Header.Referer()))
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

// handlePassthroughRequest handles requests which are supposed to be proxied to the destination host.
// If the user is authorized they are allowed to pass, otherwise they should be redirected to
// the authentication page. (/_portgate)
func (h *RequestHandler) handlePassthroughRequest(ctx *fasthttp.RequestCtx, p Path) {
	// TODO: Check authorization.
	// TODO: Check whether port is allowed to be accessed.

	// We reuse the request given to us by the user with minor changes to route it to the
	// destination host.
	ctx.Request.SetRequestURI(p.MakeUrl(h.config.targetHost))
	ctx.Request.Header.Set("Host", h.config.TargetAddress(p.DestinationIdentifier))

	// We pipe the response given to us by the destination host back to the user.
	// Since it's possible that we get a redirect, we take this into account,
	// but only allow upto 10 redirects.
	err := h.client.DoRedirects(&ctx.Request, &ctx.Response, 10)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_, _ = ctx.WriteString("An error occurred.")
	}
}

// handlePortgateRequest handles all Portgate specific request for either showing Portgate
// specific pages or handling creation of authorization tokens.
func (h *RequestHandler) handlePortgateRequest(ctx *fasthttp.RequestCtx) {
	// TODO: Implement authentication, authorization
	_, _ = ctx.WriteString("Portgate request.")
}

// handleUnknownRequest handles any request which could not be processed due to missing
// information.
func (h *RequestHandler) handleUnknownRequest(ctx *fasthttp.RequestCtx) {
	// TODO: Show error page
	ctx.SetStatusCode(http.StatusNotFound)
	_, _ = ctx.WriteString("Unknown request.")
}

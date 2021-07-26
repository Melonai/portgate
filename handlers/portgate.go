package handlers

import "github.com/valyala/fasthttp"

// handlePortgateRequest handles all Portgate specific request for either showing Portgate
// specific pages or handling creation of authorization tokens.
func (h *RequestHandler) handlePortgateRequest(ctx *fasthttp.RequestCtx) {
	// TODO: Implement authentication, authorization
	_, _ = ctx.WriteString("Portgate request.")
}

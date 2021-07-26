package handlers

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"portgate"
)

// handlePortgateRequest handles all Portgate specific request for either showing Portgate
// specific pages or handling creation of authorization tokens.
func (h *RequestHandler) handlePortgateRequest(ctx *fasthttp.RequestCtx, path portgate.Path) {
	if path.IsPortgateStaticPath() {
		h.staticHandler(ctx)
	} else {
		// TODO: Implement authentication, authorization
		h.handlePortgateIndexRequest(ctx)
	}
}

// handlePortgateIndexRequest delegates requests directed at /_portgate.
func (h *RequestHandler) handlePortgateIndexRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "GET":
		// If we received a GET request, the user wants to see either the login page,
		// or an info page if they're already authenticated.
		h.handlePortgatePageRequest(ctx)
	case "POST":
		// If we received a POST request, the user wants to get authorization to Portgate.
		h.handleAuthenticateRequest(ctx)
	}
}

// handlePortgatePageRequest renders the Portgate page with either the authentication page or
// a basic information page.
func (h *RequestHandler) handlePortgatePageRequest(ctx *fasthttp.RequestCtx) {
	// We render the page template and pass it to the user.
	ctx.Response.Header.SetContentType("text/html")
	err := h.templates.ExecuteTemplate(ctx, "authenticate.template.html", nil)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_, _ = ctx.WriteString("An error occurred.")
	}
}

func (h *RequestHandler) handleAuthenticateRequest(ctx *fasthttp.RequestCtx) {
	// TODO
}

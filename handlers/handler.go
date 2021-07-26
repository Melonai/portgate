package handlers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"portgate"
	"strings"
)

// RequestHandler keeps data relevant to the request handlers.
type RequestHandler struct {
	// Pointer to the global Portgate config, the values of which can change at runtime.
	config *portgate.Config
	// HTTP Client for requesting resources from the destination host.
	client fasthttp.Client
	// Handler for static Portgate assets.
	staticHandler fasthttp.RequestHandler
	// Templates for Portgate pages.
	templates portgate.Templates
}

// NewRequestHandler creates a new RequestHandler instance.
func NewRequestHandler(config *portgate.Config, templates portgate.Templates) RequestHandler {
	// Serves static Portgate files when called.
	fs := fasthttp.FS{
		Root: "./assets/static/",
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			return []byte(strings.TrimPrefix(string(ctx.Path()), "/_portgate/static"))
		},
		PathNotFound: nil,
	}
	staticHandler := fs.NewRequestHandler()

	return RequestHandler{
		config:        config,
		client:        fasthttp.Client{},
		staticHandler: staticHandler,
		templates:     templates,
	}
}

// HandleRequest handles all types of requests and delegates to more specific handlers.
func (h *RequestHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	path := portgate.ParsePath(string(ctx.Path()))

	if path.IsPortgatePath() {
		h.handlePortgateRequest(ctx, path)
		return
	}

	if path.DestinationIdentifier == -1 {
		// We were not given a destination.

		// Try to grab actual destination from Referer header.
		// This can help us if the user followed an absolute link on a proxied page.
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
	} else {
		// We were given a port, so we have to pass the request through to the destination host.

		h.handlePassthroughRequest(ctx, path)
	}
}

// handleUnknownRequest handles any request which could not be processed due to missing
// information.
func (h *RequestHandler) handleUnknownRequest(ctx *fasthttp.RequestCtx) {
	// TODO: Show error page
	ctx.Error("Unknown request.", http.StatusNotFound)
}

// handleUnknownRequest handles errors which occurred during a request with a generic message.
func (h *RequestHandler) handleError(ctx *fasthttp.RequestCtx) {
	// TODO: Show error page
	ctx.Error("An error occurred", http.StatusInternalServerError)
}

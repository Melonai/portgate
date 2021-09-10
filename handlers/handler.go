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
	destination := portgate.DestinationFromURL(string(ctx.Path()))

	if destination.IsPortgatePath {
		h.handlePortgateRequest(ctx, destination)
		return
	}

	if destination.Port == 0 {
		// Try to get the port from the Referer.
		destination = destination.AddReferer(string(ctx.Request.Header.Referer()))

		// Still no port?
		if destination.Port == 0 {
			h.handleUnknownRequest(ctx)
			return
		}

		portgateUrl := fmt.Sprintf("/%d%s", destination.Port, destination.Path)
		ctx.Redirect(portgateUrl, http.StatusTemporaryRedirect)
		return
	}

	h.handlePassthroughRequest(ctx, destination)
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

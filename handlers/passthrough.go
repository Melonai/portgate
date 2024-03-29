package handlers

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"portgate"
)

// handlePassthroughRequest handles requests which are supposed to be proxied to the destination host.
// If the user is authorized they are allowed to pass, otherwise they should be redirected to
// the authentication page. (/_portgate)
func (h *RequestHandler) handlePassthroughRequest(ctx *fasthttp.RequestCtx, p portgate.Destination) {
	// TODO: Check whether port is allowed to be accessed.

	// Check whether given cookie is ok, if not redirect to the authentication page.
	if !portgate.VerifyTokenFromCookie(h.config, ctx) {
		ctx.Redirect("/_portgate", http.StatusTemporaryRedirect)
		return
	}

	// We reuse the request given to us by the user with minor changes to route it to the
	// destination host.
	ctx.Request.SetRequestURI(h.config.MakeUrl(p))
	ctx.Request.Header.SetHost(h.config.TargetAddress(p.Port))

	// We pipe the response given to us by the destination host back to the user.
	// Since it's possible that we get a redirect, we take this into account,
	// but only allow upto 10 redirects.
	err := h.client.DoRedirects(&ctx.Request, &ctx.Response, 10)
	if err != nil {
		h.handleError(ctx)
	}
}

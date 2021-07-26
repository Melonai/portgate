package handlers

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"portgate"
	"time"
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
	ctx.Response.Header.SetContentType("text/html")

	var err error

	// We render the page template and pass it to the user.
	if portgate.VerifyTokenFromCookie(h.config, ctx) {
		// User is authenticated, show the information page
		err = h.templates.ExecuteTemplate(ctx, "information.template.html", nil)
	} else {
		// Show the authentication page
		err = h.templates.ExecuteTemplate(ctx, "authenticate.template.html", nil)
	}

	if err != nil {
		h.handleError(ctx)
	}
}

func (h *RequestHandler) handleAuthenticateRequest(ctx *fasthttp.RequestCtx) {

	givenKey := ctx.PostArgs().Peek("key")
	if givenKey == nil || !h.config.CheckKey(string(givenKey)) {
		ctx.Error("Wrong key.", http.StatusUnauthorized)
		return
	}

	token, err := portgate.CreateToken(h.config, string(givenKey))
	if err != nil {
		h.handleError(ctx)
	}

	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)

	cookie.SetExpire(portgate.GetExpirationDateFrom(time.Now()))
	cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	cookie.SetHTTPOnly(true)
	cookie.SetKey("_portgate_token")
	cookie.SetValue(token)

	ctx.Response.Header.SetCookie(cookie)

	// TODO: Redirect to previously request path.
	// http.StatusFound redirects a POST request to a GET request.
	ctx.Redirect("/_portgate", http.StatusFound)
}

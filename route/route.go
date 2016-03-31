package route

import (
	"net/http"

	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/controller"
	"github.com/aerth/ndjinn/route/middleware/acl"
	hr "github.com/aerth/ndjinn/route/middleware/httprouterwrapper"
	"github.com/aerth/ndjinn/route/middleware/logrequest"
	"github.com/aerth/ndjinn/route/middleware/pprofhandler"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Load the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// Load the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// Load the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.Index)))

	// Login
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
	r.GET("/logout", hr.Handler(alice.
		New().
		ThenFunc(controller.Logout)))

	// Register
	r.GET("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterGET)))
	r.POST("/register", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.RegisterPOST)))

	// Checkout
	r.GET("/checkout", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CheckoutGET)))
	r.GET("/checkout/confirm", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CheckoutConfirm)))
	r.GET("/checkout/process", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CheckoutProcess)))
	r.GET("/checkout/cancel", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CheckoutCancel)))
	r.GET("/checkout/fail", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CheckoutFail)))

	// About
	r.GET("/about", hr.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))

	// Dashboard
	r.GET("/dashboard", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.DashboardGET)))

	// Dashboard
	r.GET("/promote/{user}", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PromoteUser)))
	// Dashboard
	r.GET("/dashboard/paid", hr.Handler(alice.
		New(acl.AllowPaid).
		ThenFunc(controller.MemberDashboardGET)))

	// API Post requests
	r.POST("/api", hr.Handler(alice.
		New(acl.AllowAPI).
		ThenFunc(controller.APIPost)))

	// Status / Ping Pong x3
	r.GET("/status", hr.Handler(alice.
		New().
		ThenFunc(controller.APIStatusGet)))
	// Status / Ping Pong
	r.GET("/ping", hr.Handler(alice.
		New().
		ThenFunc(controller.APIStatusGet)))
	// Status / Ping Pong
	r.GET("/api", hr.Handler(alice.
		New(acl.AllowAPI).
		ThenFunc(controller.APIStatusGet)))

	// Enable Pprof
	r.GET("/debug/pprof/*pprof", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	cs.ExcludeRegexPaths([]string{"/api(.*)"}) // ?
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs

	// Log every request
	h = logrequest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}

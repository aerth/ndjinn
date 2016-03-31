package acl

import (
	"net/http"
"fmt"
	"github.com/aerth/ndjinn/components/session"
)

// DisallowAuth does not allow authenticated users to access the page
func DisallowAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is authenticated, don't allow them to access the page
		if sess.Values["id"] != nil {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page
func DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page
func AllowPaid(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["membership"] != 2 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// For now we are still using the session. API keys soon.
func AllowAPI(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			//printf(w, "Invalid authentication", http.StatusBadRequest)
			fmt.Println("Bunk API Request.")
			return
		}

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["membership"] != 1 {
			//http.Redirect(w, r, "/login", http.StatusFound)
			fmt.Println("Bunk API Request.")
			return
		}

		h.ServeHTTP(w, r)
	})
}

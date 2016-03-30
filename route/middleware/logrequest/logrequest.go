package logrequest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Handler will log the HTTP requests
func Handler(next http.Handler) http.Handler {
	// Sanitize user agent string for logs
	p := bluemonday.UGCPolicy()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL, p.Sanitize(r.UserAgent()))
		next.ServeHTTP(w, r)

	})
}

package logrequest

import (
	"net/http"
	"os"
	"log"
	"github.com/microcosm-cc/bluemonday"
)

// Handler will log the HTTP requests
func Handler(next http.Handler) http.Handler {
	f, e := os.OpenFile("logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if e != nil {
		log.Fatal(e)
	}
	rlog := log.New(f, "", log.LstdFlags)
	// Sanitize user agent string for logs
	p := bluemonday.UGCPolicy()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//rlog.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL, p.Sanitize(r.UserAgent()))
		rlog.Println(r.RemoteAddr, r.Method, r.URL, p.Sanitize(r.UserAgent()))
		next.ServeHTTP(w, r)

	})
}

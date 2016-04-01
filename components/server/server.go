package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	//	"github.com/fvbock/endless"
)

// Server stores the hostname and port number
type Server struct {
	Hostname  string `json:"Hostname"`  // Server name
	UseHTTP   bool   `json:"UseHTTP"`   // Listen on HTTP
	UseHTTPS  bool   `json:"UseHTTPS"`  // Listen on HTTPS
	HTTPPort  int    `json:"HTTPPort"`  // HTTP port
	HTTPSPort int    `json:"HTTPSPort"` // HTTPS port
	CertFile  string `json:"CertFile"`  // HTTPS certificate
	KeyFile   string `json:"KeyFile"`   // HTTPS private key
}

// Run starts the HTTP and/or HTTPS listener
func Run(httpHandlers http.Handler, httpsHandlers http.Handler, s Server) {
	if s.UseHTTP && s.UseHTTPS {
		go func() {
			startHTTPS(httpsHandlers, s)
		}()

		startHTTP(httpHandlers, s)
	} else if s.UseHTTP {
		startHTTP(httpHandlers, s)
	} else if s.UseHTTPS {
		startHTTPS(httpsHandlers, s)
	} else {
		log.Println("Config file does not specify a listener to start")
	}
}

// startHTTP starts the HTTP listener
func startHTTP(handlers http.Handler, s Server) {
	log.Println("Info: Listening at http://" + httpAddress(s))

	// Start the HTTP listener
	//	err := (endless.ListenAndServe(httpAddress(s), handlers))
	err := (http.ListenAndServe(httpAddress(s), handlers))
	if err != nil {
		log.Println(err)
	}
	log.Println("Server not able to listen.")
}

// startHTTPs starts the HTTPS listener
func startHTTPS(handlers http.Handler, s Server) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "- Running https://"+httpsAddress(s))

	// Start the HTTPS listener
	//	log.Fatal(endless.ListenAndServeTLS(httpsAddress(s), s.CertFile, s.KeyFile, handlers))
	log.Fatal(http.ListenAndServeTLS(httpsAddress(s), s.CertFile, s.KeyFile, handlers))
}

// httpAddress returns the HTTP address
func httpAddress(s Server) string {
	if s.Hostname == "" {
		s.Hostname = "localhost"
	}
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPPort)
}

// httpsAddress returns the HTTPS address
func httpsAddress(s Server) string {
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPSPort)
}

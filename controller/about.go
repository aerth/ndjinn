package controller

import (
	"net/http"

	"github.com/aerth/ndjinn/components/view"
)

// Displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}

package controller

import (
	"net/http"

	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
)

// Displays the default home page
func Index(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)

	if session.Values["id"] != nil {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {

		// Flash Anon Announcement

		// Display the view

		v := view.New(r)
		v.Name = "index/anon"
		session.AddFlash(view.Flash{"Welcome to " + v.GlobalSiteName + "!", view.FlashSuccess})

		v.Render(w)
		return
	}
}

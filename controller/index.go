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
		// Flash Auth Announcement

		// Display the view
		v := view.New(r)

		v.Name = "index/auth"
		v.Vars["first_name"] = session.Values["first_name"]
		v.Render(w)
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

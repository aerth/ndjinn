package controller

import (
	"net/http"

	"github.com/josephspurrier/csrfbanana"

	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
)

// Displays the About page
func DashboardGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	session := session.Instance(r)
	if session.Values["id"] != nil {
		// Flash Auth Announcement
		//	session.AddFlash(view.Flash{"Almost!", view.FlashError})

		v := view.New(r)
		v.Name = "dashboard/index"

		v.Vars["nickname"] = session.Values["nickname"]
		v.Vars["email"] = session.Values["email"]
		v.Vars["token"] = csrfbanana.Token(w, r, session)
		v.Render(w)

	} else {
		// Flash Anon Announcement
		session.AddFlash(view.Flash{"Almost!", view.FlashError})

		// Display the view
		//
		// v := view.New(r)
		// v.Name = "index/anon"
		// v.Render(w)
		//
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

// Displays the About page
func DashboardAsyncGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	session := session.Instance(r)
	if session.Values["id"] != nil {
		// Flash Auth Announcement
		//	session.AddFlash(view.Flash{"Almost!", view.FlashError})

		v := view.New(r)
		v.Name = "dashboard/async/newlisting"

		v.Vars["nickname"] = session.Values["nickname"]
		v.Vars["email"] = session.Values["email"]
		v.Vars["token"] = csrfbanana.Token(w, r, session)
		v.Render(w)

	} else {
		// Flash Anon Announcement
		session.AddFlash(view.Flash{"Anonymous!", view.FlashError})

		// Display the view
		//
		// v := view.New(r)
		// v.Name = "index/anon"
		// v.Render(w)
		//
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
}

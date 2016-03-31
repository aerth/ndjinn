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
	sess := session.Instance(r)
	if sess.Values["id"] != nil {
		// Flash Auth Announcement
		//	session.AddFlash(view.Flash{"Almost!", view.FlashError})

		v := view.New(r)
		v.Name = "dashboard/index"

		v.Vars["nickname"] = sess.Values["nickname"]
		v.Vars["email"] = sess.Values["email"]
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		sess.Save(r, w)
		v.Render(w)

	} else {
		// Flash Anon Announcement
		sess.AddFlash(view.Flash{"Almost!", view.FlashError})

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
	sess := session.Instance(r)
	if sess.Values["id"] != nil {
	if sess.Values["GlobalNotification"] != nil {
		if str, ok := sess.Values["GlobalNotification"].(string); ok {
			sess.AddFlash(view.Flash{str, view.FlashError})
			sess.Save(r, w)
			return
		}else {
			sess.Save(r, w)
			return
		}

return
}

	} else {
		// Flash Anon Announcement
		sess.AddFlash(view.Flash{"Anonymous!", view.FlashError})
		sess.Save(r, w)

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


// Displays the About page
func MemberDashboardGET(w http.ResponseWriter, r *http.Request) {
		sess := session.Instance(r)
	sess.Save(r, w)
	// Display the view
	v := view.New(r)
	v.Name = "dashboard/members/index"
	v.Render(w)
}

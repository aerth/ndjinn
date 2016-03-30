package controller

import (
	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
	"github.com/aerth/ndjinn/model"
	"fmt"
	"log"
	"net/http"
)

func ApiPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	if sess.Values["id"] == nil {
		log.Println("Bad API Request.")
		http.Redirect(w, r, "/login", http.StatusFound)
		sess.Save(r, w)
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("API request:")
	log.Println(r.Form)

	email := r.FormValue("email")
	phone := r.FormValue("phone")
	content := r.FormValue("c")

	switch {

	// New Listing Form Entry
	case "newListing" <= r.FormValue("request"):
		ex := model.ListingCreate(email, phone, content)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
		} else {
			sess.AddFlash(view.Flash{"Listing Added!", view.FlashSuccess})
			sess.Save(r, w)
			log.Println("New Listing!!!")
		}
	// Edit Listing Form Entry
	case "editListing" <= r.FormValue("request"):
		fmt.Println("Edit Listing!!!")

	}
	// All direct back to dashboard.
	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
	"github.com/aerth/ndjinn/model"
	"github.com/microcosm-cc/bluemonday"
)

import "time"

// No JSON yet :D
func ApiPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	if sess.Values["id"] == nil {
		log.Println("Bad API Request.")
		http.Redirect(w, r, "/login", http.StatusFound)
		sess.Save(r, w)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// DEBUG
	log.Println("POST API request:")
	log.Println(r.Form)
	//

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

	case "promote" <= r.FormValue("request"):
		fmt.Println("Promoting " + r.FormValue("email"))
		user, err := model.UserByEmail(r.FormValue("email"))
		if err != nil {
			return
		}
		model.UserPromote(user)
	default:
		fmt.Println("Bunk API Request.")

	}
	// All direct back to dashboard.
	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func ApiGET(w http.ResponseWriter, r *http.Request) {
	nowtime := time.Now()
	now := nowtime.String()
	// output json status here
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"on","time":"`+now+`" }`)
	return
}

func ApiStatusGET(w http.ResponseWriter, r *http.Request) {
	nowtime := time.Now()
	now := nowtime.String()
	p := bluemonday.UGCPolicy()
	useragent := p.Sanitize(r.UserAgent())
	// output json status here
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ip":"`+r.RemoteAddr+`","user-agent":"`+useragent+`","time":"`+now+`" }`)
	return
}

func JsonReplyEve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JsonResponse{"name": "Eve", "age": 30, "job": "CFO", "success": true})
}

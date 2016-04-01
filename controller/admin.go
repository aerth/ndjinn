package controller

import (
	"github.com/aerth/ndjinn/components/admin"
	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
	"github.com/aerth/ndjinn/model"
	//	"github.com/aerth/ndjinn/model"

	"log"
	"net/http"

	"github.com/josephspurrier/csrfbanana"
)

var (
	a          admin.AdminInfo
	newAccount bool = true
)

func AdminGet(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	sess.AddFlash(view.Flash{"Admin!", view.FlashError})
	sess.Save(r, w)
	// Display the view
	v := view.New(r)
	v.Name = "dashboard/index"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"nickname", "email"}, r.Form, v.Vars)

	v.Render(w)
}
func init() {

}

func AdminPromoteUser(w http.ResponseWriter, r *http.Request) {
	//now := time.Now()
	sess := session.Instance(r)

	email := sess.Values["email"]
	user, err := model.UserByEmail(email.(string))
	if err != nil {
		log.Println("line 297")
	}
	//	var membership int = 4
	ex := model.UserPromote(user)
	// Will only error if there is a problem with the query
	if ex != nil {
	}

	sess.Save(r, w)
	v := view.New(r)
	v.Name = "dashboard/members"
	v.Render(w)

}

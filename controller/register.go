package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aerth/ndjinn/components/passhash"
	"github.com/aerth/ndjinn/components/recaptcha"
	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
	"github.com/aerth/ndjinn/model"

	"github.com/josephspurrier/csrfbanana"
)

func RegisterGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "register/register"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"NickName", "Email"}, r.Form, v.Vars)
	v.Render(w)
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		sess.AddFlash(view.Flash{"Please try again later.", view.FlashError})
		sess.Save(r, w)
		log.Println("Brute force register prevented")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"NickName", "Email", "GoodPassword"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}

	// Get form values
	nickName := r.FormValue("NickName")
	membershipLevel := strconv.Itoa(0)
	email := r.FormValue("Email")
	password, errp := passhash.HashString(r.FormValue("GoodPassword"))

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later. Error Code: P67", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Get database result
	_, err := model.UserByEmail(email)

	if err == model.ErrNoResult { // If success (no user exists with that email)
		ex := model.UserCreate(nickName, membershipLevel, email, password)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.  Error Code: Q81", view.FlashError})
			sess.Save(r, w)
		} else {
			sess.AddFlash(view.Flash{"Account created successfully for: " + email, view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later. Error Code: 91", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		sess.AddFlash(view.Flash{"Account already exists for: " + email, view.FlashError})
		sess.Save(r, w)
	}

	// Display the page
	RegisterGET(w, r)
}

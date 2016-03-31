package controller

import (
	"github.com/aerth/ndjinn/components/checkout"
	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/model"
	"github.com/aerth/ndjinn/components/view"
//	"github.com/aerth/ndjinn/model"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/josephspurrier/csrfbanana"

	paypalsdk "github.com/logpacker/PayPal-Go-SDK"
)

var (
	z          checkout.CheckoutInfo
	newAccount bool = true
)

func CheckoutGETOLD(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	sess.AddFlash(view.Flash{"Almost Done! Just need 10 bucks!", view.FlashError})
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

func CheckoutGET(w http.ResponseWriter, r *http.Request) {
	//var buf bytes.Buffer
	now := time.Now()

	nowtime := now.Format("Mon, 2 Jan 2006 15:04:05 -0700")
		v := view.New(r)
	sess := session.Instance(r)
	sess.Save(r, w) // Ensure you save the session after making a change to it

	//	sess := session.Instance(r)
	//	email := sess.Values["email"]
	//membership := sess.Values["membership"]
	//email := sess.Get("email")
	//if email != "" {fmt.Println(email)}

	z = checkout.ReadConfig()
	log.Println(z.PayPalC)
	log.Println(z.PayPalK)
	var c, err = paypalsdk.NewClient(z.PayPalC, z.PayPalK, paypalsdk.APIBaseSandBox)
	if err == nil {
		log.Println("ClientID: " + c.ClientID)
		log.Println("APIBase: " + c.APIBase)
	} else {
		log.Println("ERROR: " + err.Error())
		return
	}
	log.Println("Getting new AccessToken")
	token, err := c.GetAccessToken()
	if err == nil {
		log.Println("AccessToken: " + token.Token)

	} else {
		v.Name = "dashboard/index"

		sess.AddFlash(view.Flash{"Internal Error. Please contact support with this timestamp:", view.FlashNotice})
		sess.AddFlash(view.Flash{nowtime + r.RemoteAddr + r.RequestURI, view.FlashNotice})
		sess.Save(r, w) // Ensure you save the session after making a change to it
		fmt.Println("ERROR: " + err.Error())
		v.Render(w)
		return
	}

	amount := paypalsdk.Amount{
		Total:    "10.00",
		Currency: "USD",
	}
	redirectURI := v.BaseURI + "checkout/confirm"
	cancelURI := v.BaseURI + "checkout/cancel"

	description := "MembershipLevel (one year) for " + nowtime
	log.Println(description)
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
	/*
	   payment, err := c.GetPayment(paymentResult.ID)
	   payments, err := c.GetPayments()
	   if err == nil {
	     fmt.Println("DEBUG: PaymentsCount=" + strconv.Itoa(len(payments)))
	   } else {
	     fmt.Println("ERROR: " + err.Error())
	   }
	   fmt.Println("Payment ID:")
	   fmt.Println(payment.ID)
	   fmt.Println("Payment:")
	   fmt.Println(payment)
	   fmt.Println("Payments:")
	   fmt.Println(payments)
	*/

	//http.Redirect(w, r, paymentResult.Links[1].Href, 302)
	if len(paymentResult.Links) > 2 {
	fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\""+paymentResult.Links[1].Href+"\">Click here to use PayPal ($ 10)</a></html>")
	log.Println("Redirecting " + r.RemoteAddr + paymentResult.ID + paymentResult.Links[1].Href)
	}
	return
}

func CheckoutConfirm(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	nowtime := now.Format("Mon, 2 Jan 2006 15:04:05 -0700")
	sess := session.Instance(r)
	sess.Save(r, w) // Ensure you save the session after making a change to it
	v := view.New(r)
	// m, e := url.ParseQuery(r.URL.RawQuery)
	// if e != nil {
	// 	log.Println(e)
	// 	return
	// }
	log.Println("Returned from PayPal.")
	//log.Println(m["id"])
	a, e := url.ParseQuery(r.URL.RawQuery)
	if e != nil {
		log.Println(e)
		return
	}
	//paymentid := strings.Join(m["id"][:0], "")
	//payerid := strings.Join(m["u"][:0], "")
	paymentid := a.Get("id")
	payerid := a.Get("u")
	if payerid == "" || paymentid == "" {
		paymentid = a.Get("paymentId")
		payerid = a.Get("PayerID")
	}
	log.Println("Payment ID: " + paymentid)
	log.Println("Payer ID: " + payerid)
	if payerid == "" || paymentid == "" {
		v.Name = "dashboard/index"

		sess.AddFlash(view.Flash{"Internal Error. Your payment was not processed.", view.FlashNotice})

		sess.Save(r, w) // Ensure you save the session after making a change to it
		fmt.Println("ERROR. Blank PaymentID or PayerID")
		v.Render(w)

		return
	}
	z = checkout.ReadConfig()

	var c, err = paypalsdk.NewClient(z.PayPalC, z.PayPalK, paypalsdk.APIBaseSandBox)
	if err == nil {
		log.Println("ClientID: " + c.ClientID)
		log.Println("APIBase: " + c.APIBase)
	} else {
		v.Name = "dashboard/index"
		sess := session.Instance(r)
		sess.AddFlash(view.Flash{"Internal Error. Please contact support with this timestamp:", view.FlashNotice})
		sess.AddFlash(view.Flash{nowtime + r.RemoteAddr + r.RequestURI, view.FlashNotice})
		sess.Save(r, w) // Ensure you save the session after making a change to it
		fmt.Println("ERROR: " + err.Error())
		v.Render(w)

		return
	}

	log.Println("Getting new AccessToken")
	token, err := c.GetAccessToken()
	if err == nil {
		log.Println("AccessToken: " + token.Token)

	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	payment, err := c.GetPayment(paymentid)
	if err != nil {
		log.Println("Paypal Error")
		log.Println(err.Error())
		http.Redirect(w, r, "/checkout/fail?id="+paymentid, 302)
	} else {
		log.Println("Pre-Confirm:")
		log.Println(payment.Intent, payment.Payer.PayerInfo.FirstName, payment.Payer.PayerInfo.LastName)
		log.Println(payment.Payer.PayerInfo.Email, payment.Payer.PayerInfo.Phone, payment.Payer.PayerInfo.PayerID)
		transaction := payment.Transactions
		fmt.Fprintf(w, "<!DOCTYPE html><html><a href=\"/checkout/process?id="+payment.ID+"&u="+payment.Payer.PayerInfo.PayerID+"\">Click here to Confirm Payment of $%s</a></html>", transaction[0].Amount.Total)
	}

}
func CheckoutCancel(w http.ResponseWriter, r *http.Request) {

	m, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(m)

}
func CheckoutFail(w http.ResponseWriter, r *http.Request) {
	log.Println("Failed.")
	c, err := paypalsdk.NewClient(z.PayPalC, z.PayPalK, paypalsdk.APIBaseSandBox)
	if err != nil {
		log.Println(err)
	}
	m, _ := url.ParseQuery(r.URL.RawQuery)
	paymentid := m["id"][0]
	log.Println(paymentid)
	payment, err := c.GetPayment(paymentid)
	if payment != nil {
		log.Println(payment)
	}
	if err != nil {
		log.Println(err.Error())
	}
	if payment != nil {
		log.Println(payment.ID)
	}

	http.Redirect(w, r, "/", 302)

}

func CheckoutProcess(w http.ResponseWriter, r *http.Request) {
	//now := time.Now()
	//nowtime := now.Format("Mon, 2 Jan 2006 15:04:05 -0700")
	v := view.New(r)
	log.Println("Payment Confirmed by User. Sending to Paypal.")
	c, err := paypalsdk.NewClient(z.PayPalC, z.PayPalK, paypalsdk.APIBaseSandBox)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	m, _ := url.ParseQuery(r.URL.RawQuery)
	paymentid := m["id"][0]
	payerid := m["u"][0]
	//payment, err := c.GetPayment(paymentid)
	log.Println("Payment ID: " + paymentid)
	log.Println("Payer ID: " + payerid)
	c.GetAccessToken()
	response, err := (c.ExecuteApprovedPayment(paymentid, payerid))
	log.Println(response)
	if err != nil {
		log.Println("Error:")
		log.Println(err.Error())
		return
	}
	if response.State != "" {
		log.Println("Approval State: " + response.State)
	}

	if response.State != "approved" {
		log.Println("Something wrong.")
		return
	}
	v.Name = "dashboard/index"
	sess := session.Instance(r)
	//email := sess.Values["email"]

	//var user = &model.User{}
	//
	// if user, ok := email.(*&model.User); !ok {
	// 	log.Println("line 262")
	// 	return
	// }

	//userpromote, err := model.UserByEmail(user.Email)
	// log.Println(user)
	// log.Println(user.Email)
	// log.Println(userpromote)
	// model.UserPromote(userpromote)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	//log.Println(sess)
	sess.AddFlash(view.Flash{"Thank you for being a valued member.", view.FlashNotice})
	sess.Save(r, w) // Ensure you save the session after making a change to it

	fmt.Println("GOT PAID *********")

	v.Render(w)

	return
}

func PromoteUser(w http.ResponseWriter, r *http.Request) {
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
		if ex != nil {}

			sess.Save(r, w)
				v := view.New(r)
		v.Name = "dashboard/members"
	v.Render(w)

}

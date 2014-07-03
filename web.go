package main

import (
	"github.com/joho/godotenv"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	//base URL of this application
	BASE_URL string

	//Outgoing Caller ID you have previously validated with Twilio
	CALLER_ID string

	client Client
)

func index(c web.C, w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	fmt.Fprintf(w, `<h1>Twilio phone reminder demo</h1>
<h2 style="color: #ff0000">%s</h2>
<h3>Enter your phone number to receive an automated reminder</h3>
<form action="/makecall" method="post">
    <input type="text" name="number" />
    <input type="submit" value="Call me!">
</form>`, msg)
}

func makecall(c web.C, w http.ResponseWriter, r *http.Request) {
	to := r.PostFormValue("number")
	if len(to) == 0 {
		fmt.Fprintln(w, "Number required")
		return
	}
	attrs := url.Values{
		"From":           {CALLER_ID},
		"Url":            {BASE_URL + "/reminder"},
		"To":             {to},
		"StatusCallback": {BASE_URL + "/statusCallback"},
	}
	resp, err := client.CreateCall(attrs)
	if err != nil {
		fmt.Fprintln(w, "Bang! "+err.Error())
		return
	}
	fmt.Fprintf(w, "Calling %s...", to)
	fmt.Print("Twilio said", string(resp))
}

func reminder(c web.C, w http.ResponseWriter, r *http.Request) {
	postTo := BASE_URL + "/directions"
	fmt.Fprint(w, ReminderXml(postTo))
}

func directions(c web.C, w http.ResponseWriter, r *http.Request) {
	digits := r.FormValue("Digits")
	fmt.Printf("Got digit %s from Twilio", digits)

	if digits == "3" {
		http.Redirect(w, r, "/goodbye", http.StatusFound)
		return
	}
	if len(digits) == 0 || digits != "2" {
		http.Redirect(w, r, "/reminder", http.StatusFound)
		return
	}
	redirectTo := BASE_URL + "/reminder"
	fmt.Fprint(w, DiretionsXml(redirectTo))
}

func statusCallback(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("Callback info %v\n", r.Form)
	fmt.Fprint(w, "ok")
}

func goodbye(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, GoodbyeXml())
}

func main() {
	//your Twilio authentication credentials
	godotenv.Load()
	ACCOUNT_SID := os.Getenv("ACCOUNT_SID")
	ACCOUNT_TOKEN := os.Getenv("ACCOUNT_TOKEN")
	BASE_URL = os.Getenv("BASE_URL")
	CALLER_ID = os.Getenv("CALLER_ID")
	client = *NewClient(ACCOUNT_SID, ACCOUNT_TOKEN)
	goji.Get("/", index)
	goji.Post("/makecall", makecall)
	goji.Get("/reminder", reminder)
	goji.Post("/reminder", reminder)
	goji.Get("/directions", directions)
	goji.Post("/directions", directions)
	goji.Get("/goodbye", goodbye)
	goji.Post("/goodbye", goodbye)

	goji.Get("/statusCallback", statusCallback)
	goji.Post("/statusCallback", statusCallback)

	goji.Serve()
}

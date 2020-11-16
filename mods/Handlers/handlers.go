package handlers

import (
	db "DB"
	pages "Pages"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var (
	//SessionStore is the session store from main
	SessionStore *sessions.CookieStore
	//SessionName is the session name from main
	SessionName string
)

//IndexHandler handles the "/index/" request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.Index{}
	session, err := SessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
		return
	}
	username, found := session.Values["username"]
	if !found || username == "" {
		data.Title = "Please log in or register"
		data.LoggedIn = false
	} else {
		data.Title = "Welcome " + fmt.Sprintf("%v", username)
		data.LoggedIn = true
		data.Sheets = db.GetSheets(fmt.Sprintf("%v", username))
	}
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, data)
}

//LoginPageHandler loads the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

//LoginHandler handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if db.CheckUser(username, password) {
		session, err := SessionStore.Get(r, SessionName)
		if err != nil {
			handleSessionError(w, err)
			return
		}
		session.Values["username"] = username
		if err := session.Save(r, w); err != nil {
			handleSessionError(w, err)
			return
		}
	}
	http.Redirect(w, r, "/index/", 303)
}

//LogoutHandler logs the user out and expires the session cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := SessionStore.Get(r, "cookie-name")
	if err != nil {
		handleSessionError(w, err)
		return
	}
	fmt.Println(session.Values["username"])
	session.Values["username"] = ""
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/index/", 303)
}

//HandleSessionError is called when an error occurs with the session.
func handleSessionError(w http.ResponseWriter, err error) {
	http.Error(w, "Application Error", http.StatusInternalServerError)
	fmt.Println("Error happened")
}

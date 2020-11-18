package main

import (
	db "DB"
	pages "Pages"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

//SessionStore sotres the sessions
var sessionStore *sessions.CookieStore

//SessionName is the name of the session
var sessionName = "Sheet-Session"

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/loginpage/", loginPageHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/sheet/", sheetHandler)
	log.Printf("Listening on %s...\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:   "session",
			Value:  encoded,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func getUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func determineEncryptionKey() ([]byte, error) {
	sek := os.Getenv("SESSION_ENCRYPTION_KEY")
	lek := len(sek)
	switch {
	case lek >= 0 && lek < 16, lek > 16 && lek < 24, lek > 24 && lek < 32:
		return nil, errors.Errorf("SESSION_ENCRYPTION_KEY needs to be either 16, 24 or 32 characters long or longer, was: %d", lek)
	case lek == 16, lek == 24, lek == 32:
		return []byte(sek), nil
	case lek > 32:
		return []byte(sek[0:32]), nil
	default:
		return nil, errors.New("invalid SESSION_ENCRYPTION_KEY: " + sek)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.Index{}
	username := getUserName(r)
	if username == "" {
		data.Title = "Please log in or register"
		data.LoggedIn = false
	} else {
		data.Title = "Welcome " + username
		data.LoggedIn = true
		data.Sheets = db.GetSheets(username)
	}
	pageData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, string(pageData))
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/login.html")
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if db.CheckUser(username, password) {
		setSession(username, w)
	}
	http.Redirect(w, r, "/index/", 303)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/index/", 303)
}

func sheetHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	page := pages.SheetPage{}
	if username == "" {
		page.LoggedIn = false
	} else {
		page.LoggedIn = true
		page.CharacterSheet = db.GetSheet(username, r.FormValue("sheet"))
	}
	pageData, err := json.Marshal(page)
	if err != nil {
		panic(err)
	}
	t, _ := template.ParseFiles("./templates/sheet.html")
	t.Execute(w, string(pageData))
}

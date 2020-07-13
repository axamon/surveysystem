package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"survey/ldaplogin"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("testdisurvey")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "surveyCTIO")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "vietato", http.StatusForbidden)

		return
	}

	// Print secret message
	fmt.Fprintln(w, "Questo lo leggi solo se ti sei autenticato")
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "surveyCTIO")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "surveyCTIO")

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/login.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		username := r.FormValue("username")
		password := r.FormValue("password")
		ok, err := ldaplogin.IsOK(username, password)
		if err != nil {
			http.Redirect(w, r, "/login", 301)
			return
		}

		// Set user as authenticated
		if ok {
			session.Values["authenticated"] = true
		}
		session.Save(r, w)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {

	var address = flag.String("addr", ":8080", "Server address")
	flag.Parse()

	http.HandleFunc("/", login)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static", fs)

	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"survey/ldaplogin"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	token = make([]byte, 32)
	key   = []byte(fmt.Sprint(rand.Read(token)))
	store = sessions.NewCookieStore(key)
)

func main() {
	var address = flag.String("addr", ":8080", "Server address")
	flag.Parse()

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/survey", survey)

	err := http.ListenAndServe(*address, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "surveyCTIO")
	if err != nil {
		log.Println(err)
	}
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Logout effettuato per %s", session.Values["utente"].(string))
	// w.Write([]byte("Logout effettuato\n"))
	// time.Sleep(2 * time.Second)
	r.Method = "GET"
	http.Redirect(w, r, "/", http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "surveyCTIO")
	if err != nil {
		log.Println(err)
	}

	switch r.Method {
	case "GET":
		survey(w, r)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)

		matricola := r.FormValue("username")
		password := r.FormValue("password")

		ok, nomeCognome, err := ldaplogin.IsOK(matricola, password)
		if err != nil {
			//http.Redirect(w, r, "/login", 301)
			log.Println(err)
		}
		// ripulisci passoword
		password = "******"
		// Set user as authenticated
		if ok {
			session.Values["authenticated"] = true
			session.Values["matricola"] = matricola
			session.Values["utente"] = nomeCognome
		}
		// Allows Admin to enter without LDAP authentication
		if matricola == "Admin" {
			session.Values["authenticated"] = true
			session.Values["matricola"] = "Admin"
			session.Values["utente"] = "Admin"

		}
		session.Save(r, w)
		r.Method = "GET"
		survey(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "static") {
			return
		}
		session, _ := store.Get(r, "surveyCTIO")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "vietato", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

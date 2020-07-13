package main

import (
    "fmt"
	"survey/ldaplogin"
	"log"
	"net/http"

    "github.com/gorilla/sessions"
)

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Print secret message
    fmt.Fprintln(w, "Questo lo leggi solo se ti sei autenticato")
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	
    // if r.URL.Path != "/" {
    //     http.Error(w, "404 not found.", http.StatusNotFound)
    //     return
    // }
 
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
		ok, err := ldaplogin.IsOK(username,password)
		if err != nil {
			fmt.Fprintf(w, "Credenziali errate")
			return
		}

		// Set user as authenticated
		if ok {
		session.Values["authenticated"] = true
		}
		session.Save(r, w)
        http.Redirect(w, r, "/secret", http.StatusMovedPermanently)
    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }
}
 

func main() {

	http.HandleFunc("/", login)


	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static", fs)

    http.HandleFunc("/secret", secret)
    http.HandleFunc("/login", login)
    http.HandleFunc("/logout", logout)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
	  log.Fatal(err)
	}
  
}
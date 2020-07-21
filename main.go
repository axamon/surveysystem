package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"net/http"
	"strings"
	"html/template"

	"github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var templates map[string]*template.Template

//Compile view templates
func init() {
	if templates == nil {
		templates=make(map[string]*template.Template)
	}
	templates["index"]=template.Must(template.ParseFiles("templates/index.gohtml", "templates/footer.gohtml"))
	templates["login"]=template.Must(template.ParseFiles("templates/index.gohtml", "templates/footer.gohtml"))
	templates["logout"]=template.Must(template.ParseFiles("templates/logout2.gohtml", "templates/footer.gohtml"))
	templates["survey"]=template.Must(template.ParseFiles("templates/survey.gohtml", "templates/footer.gohtml"))
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	token = make([]byte, 32)
	key   = []byte(fmt.Sprint(rand.Read(token)))
	store = sessions.NewCookieStore(key)
)

func main() {
	var address = flag.String("addr", ":8080", "Server address")
	flag.Parse()

	r := http.NewServeMux()
	//r := mux.NewRouter()
	
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))
	//r.Handle("/static/", fs)
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/survey", survey)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	err := http.ListenAndServe(*address, loggedRouter)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = template.Must(template.ParseFiles("templates/index.gohtml", "templates/footer.gohtml"))

	err := indexTmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
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

package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

var templates map[string]*template.Template

// Compila i templates e li inserisce nella mappa templates.
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("templates/index.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
	templates["login"] = template.Must(template.ParseFiles("templates/index.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
	templates["logout"] = template.Must(template.ParseFiles("templates/logout.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
	templates["survey"] = template.Must(template.ParseFiles("templates/survey.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
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

	var url = "http://127.0.0.1" + *address
	var err error

	switch runtime.GOOS {

	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/survey", survey)

	loggedRouter := handlers.LoggingHandler(os.Stdout, middleware(r))
	err = http.ListenAndServe(*address, loggedRouter)
	if err != nil {
		log.Fatal(err)
	}
}

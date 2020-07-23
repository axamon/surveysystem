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

	"github.com/gobuffalo/packr"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

var templates *template.Template
var templatefs = packr.NewBox("./templates")

// Compila i templates e li inserisce nella mappa templates.
func init() {
	templates = template.Must(template.ParseGlob("./templates/*.gohtml"))
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	token = make([]byte, 32)
	key   = []byte(fmt.Sprint(rand.Read(token)))
	store = sessions.NewCookieStore(key)
)

var staticfs = packr.NewBox("./static")

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

	// fs := http.FileServer(http.Dir("./static"))
	fs := http.FileServer(staticfs)
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

// ParseInternalTemplate parses templates that are not stored on filesyste.
func ParseInternalTemplate(t *template.Template, data ...InternalTemplate) (*template.Template, error) {

	for _, internalTemp := range data {
		var tmpl *template.Template
		if t == nil {
			t = template.New(internalTemp.Name)
		}
		var name string
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(internalTemp.Name)
		}
		_, err := tmpl.Parse(internalTemp.Text)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

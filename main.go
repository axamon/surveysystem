package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

var templates *template.Template
<<<<<<< HEAD
var templatefs = packr.NewBox("./templates")

// Compila i templates e li inserisce nella mappa templates.
func init() {
	templates = template.Must(template.ParseGlob("./templates/*.gohtml"))
}
=======
var store *sessions.CookieStore

//go:generate go-bindata -fs static templates
>>>>>>> b2ff8f323abf46180ced9a5358a03b91dfc82293

func init() {
	// Compila i templates e li inserisce nella mappa templates.
	go istantiateInternalTemplates()
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	token := make([]byte, 32)
	key := []byte(fmt.Sprint(rand.Read(token)))
	store = sessions.NewCookieStore(key)
}

var staticfs = packr.NewBox("./static")

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Fatal("Timeout")
	}()

	var (
		err              error
		r                *http.ServeMux
		fs, loggedRouter http.Handler
		url              string
		address          *string
	)
	address = flag.String("addr", ":8080", "Server address")
	flag.Parse()

	url = "http://127.0.0.1" + *address
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

	r = http.NewServeMux()

<<<<<<< HEAD
	// fs := http.FileServer(http.Dir("./static"))
	fs := http.FileServer(staticfs)
	r.Handle("/static/", http.StripPrefix("/static/", fs))
=======
	fs = http.FileServer(AssetFile()) // http.Dir("./static"))
	r.Handle("/static/", fs)          // http.StripPrefix("/static/", fs))
>>>>>>> b2ff8f323abf46180ced9a5358a03b91dfc82293
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/survey", survey)
	r.HandleFunc("/exit", exit)

	loggedRouter = handlers.LoggingHandler(os.Stdout, middleware(r))
	err = http.ListenAndServe(*address, loggedRouter)
	if err != nil {
		log.Fatal(err)
	}

}

func istantiateInternalTemplates() {
	t := template.New("surveysystem")

	for _, asset := range AssetNames() {
		if !strings.Contains(asset, "gohtml") {
			continue
		}
		details, err := AssetInfo(asset)
		if err != nil {
			log.Println(err)
		}
		bytes, err := Asset(details.Name())
		if err != nil {
			log.Println(err)
		}
		tmpl := t.New(strings.Split(details.Name(), "/")[1])
		_, err = tmpl.Parse(string(bytes))
		if err != nil {
			log.Println(err)
		}
	}

	templates = t
	return
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

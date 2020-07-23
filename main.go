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

var templates *template.Template

// Compila i templates e li inserisce nella mappa templates.
func init() {
	templates = template.Must(template.ParseFiles(_filePath("templates", "index.gohtml"),
		_filePath("templates", "header.gohtml"),
		_filePath("templates", "footer.gohtml"),
		_filePath("templates", "survey.gohtml"))) // "./templates/*.gohtml"))
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

	for i, a := range AssetNames() {
		fmt.Println(i, a)
	}

	a, err := templatesLogoutGohtml()

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

	fs := http.FileServer(AssetFile()) // http.Dir("./static"))
	r.Handle("/static/", fs)           // http.StripPrefix("/static/", fs))
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

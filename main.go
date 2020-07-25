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

//go:generate go-bindata -fs static templates

// Compila i templates e li inserisce nella mappa templates.
func init() {
	o, _ := templatesIndexGohtml()
	i := InternalTemplate{Name: "index.gohtml", Text: string(o.bytes)}
	o, _ = templatesHeaderGohtml()
	h := InternalTemplate{Name: "header.gohtml", Text: string(o.bytes)}
	o, _ = templatesFooterGohtml()
	f := InternalTemplate{Name: "footer.gohtml", Text: string(o.bytes)}
	o, _ = templatesGrazieGohtml()
	g := InternalTemplate{Name: "grazie.gohtml", Text: string(o.bytes)}
	o, _ = templatesErrorGohtml()
	e := InternalTemplate{Name: "error.gohtml", Text: string(o.bytes)}
	o, _ = templatesSurveyGohtml()
	s := InternalTemplate{Name: "survey.gohtml", Text: string(o.bytes)}
	o, _ = templatesLogoutGohtml()
	l := InternalTemplate{Name: "logout.gohtml", Text: string(o.bytes)}
	o, _ = templatesLoginGohtml()
	li := InternalTemplate{Name: "login.gohtml", Text: string(o.bytes)}

	t := template.New("surveysystem")
	templates = template.Must(ParseInternalTemplate(t, i, h, f, g, e, s, l, li))

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

	// t := template.New("surveysystem")
	// var internals []InternalTemplate
	// for i, name := range AssetNames() {

	// 	if strings.Contains(name, "gohtml") {
	// 		var internal InternalTemplate
	// 		body, err := Asset(name)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		internal.Name = name
	// 		internal.Text = string(body)
	// 		internals = append(internals, internal)
	// 		fmt.Println(i, name, string(body))

	// 	}
	// }
	// templates = template.Must(ParseInternalTemplate(t, internals))

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

// ParseInternalTemplate parsa i template senza recuperarli da file.
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

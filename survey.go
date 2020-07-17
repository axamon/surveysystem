package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func survey(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		session, _ := store.Get(r, "surveyCTIO")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Autenticazione errata o assente.", http.StatusForbidden)
			return
		}
		// Se autenticato...
		tmpl := template.Must(template.ParseFiles("templates/survey.gohtml"))

		data, err := ioutil.ReadFile("surveys/primo.xml")
		if err != nil {
			log.Println(err)
		}
		note := &Survey{}
		err = xml.Unmarshal([]byte(data), &note)
		if err != nil {
			log.Println(err)
		}
		note.Utente = strings.Split(session.Values["utente"].(string), " ")[0] // Aggiunge nome utente
		note.Matricola = session.Values["matricola"].(string)
		inizio, _ := time.Parse("20060102", note.Inizio)
		fine, _ := time.Parse("20060102", note.Fine)
		note.Inizio = inizio.Format("2006-01-02")
		note.Fine = fine.Format("2006-01-02")
		err = tmpl.Execute(w, note)
		if err != nil {
			log.Println(err)
		}

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		for k, v := range r.Form {
			fmt.Printf("%s = %s\n", k, v)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

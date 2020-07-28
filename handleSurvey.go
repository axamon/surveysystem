package main

import (
	"encoding/xml"
	"fmt"
<<<<<<< HEAD
	"io/ioutil"
=======
>>>>>>> b2ff8f323abf46180ced9a5358a03b91dfc82293
	"log"
	"net/http"
	"strings"
	"time"
)

func survey(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		session, _ := store.Get(r, "surveyCTIO")
		o, err := staticPrimoXml()
		if err != nil {
			log.Println(err)
		}
		data := o.bytes // ioutil.ReadFile("surveys/primo.xml")

		note := &Survey2{}
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

		// Serve template
		err = templates.ExecuteTemplate(w, "survey.gohtml", note)
		if err != nil {
			log.Println(err)
		}

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		err := writeToCSV(r.Form)
		if err != nil {
			log.Println(err)
		}

		err = templates.ExecuteTemplate(w, "grazie.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

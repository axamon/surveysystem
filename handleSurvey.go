package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func survey(w http.ResponseWriter, r *http.Request) {
	var grazieTmpl = template.Must(template.ParseFiles("templates/grazie.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
	var surveyTmpl = template.Must(template.ParseFiles("templates/survey.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))

	switch r.Method {

	case "GET":
		session, _ := store.Get(r, "surveyCTIO")
		data, err := ioutil.ReadFile("surveys/primo.xml")
		if err != nil {
			log.Println(err)
		}
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

		// Crea file csv.
		var fileCSV = "surveyID" + note.ID + ".csv"
		err = createFileCsv(fileCSV, len(note.Domande.Domanda))
		if err != nil {
			log.Printf("csv crearion in error: impossibile creare file csv: %s\n", fileCSV)
		}

		// Serve template
		err = surveyTmpl.Execute(w, note)
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

		err = grazieTmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

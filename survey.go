package main

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var csvlock sync.RWMutex

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
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		err := writeToCSV(r.Form)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Grazie per aver partecipato\n")

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func writeToCSV(data map[string][]string) error {
	csvlock.Lock()
	defer csvlock.Unlock()
	var fileCSV = "surveyID" + strings.Join(data["surveyID"], "") + ".csv"

	err := createFileCsv(fileCSV)

	if err != nil {
		return errors.New("csv crearion in error: impossibile creare file csv")
	}

	f, err := os.OpenFile(fileCSV, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	defer csvwriter.Flush()

	var record []string
	for i := 1; i <= len(data); i++ {
		if v, ok := data[strconv.Itoa(i)]; ok {
			record = append(record, strings.Join(v, ","))
		}
	}
	matricola := "\"" + strings.Join(data["matricola"], "") + "\""
	record = append(record, time.Now().Format("20060102T15:04"))
	record = append(record, matricola)
	fmt.Println(len(record))
	for i, r := range record {
		fmt.Println(i, r)
	}

	for k, v := range data {
		fmt.Println(k, v)
	}
	csvwriter.Comma = ';'

	err = csvwriter.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func createFileCsv(path string) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		file.Close()
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		defer f.Close()
		csvheader := []string{"#timestamp", "filename", "id", "titolo", "V", "N", "Guiraud", "hapax%", "italianwords%", "sigle"}
		w := csv.NewWriter(f)
		w.Write(csvheader)
		w.Flush()

		fmt.Println("Created file: ", path)
	}

	return nil
}

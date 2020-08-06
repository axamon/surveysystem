package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func survey(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		session, _ := store.Get(r, "surveyCTIO")

		uri := r.RequestURI

		sheetID := strings.Split(uri, "/")[2]

		var data []byte
		var done = make(chan struct{}, 1)
		go func() {
			data = readSheet(sheetID)
			done <- struct{}{}
		}()
		var newsurvey = new(Survey3)
		newsurvey.Utente = strings.Split(session.Values["utente"].(string), " ")[0]
		newsurvey.Matricola = session.Values["matricola"].(string)
		newsurvey.Department = session.Values["department"].(string)

		var m = make(map[string][]string)
		<-done
		s := bufio.NewScanner(strings.NewReader(string(data)))
		for s.Scan() {
			list := strings.Split(s.Text(), ",")
			m[list[0]] = list[1:]
		}
		newsurvey.TimestampInizio = time.Now().Format("02/01/2006 15:04:05")
		newsurvey.Titolo = m["Nome Survey"][0]
		newsurvey.ID = m["SurveyID"][0]
		newsurvey.Video = m["Video"][0]
		newsurvey.Inizio = m["Inizio"][0]
		newsurvey.Fine = m["Fine"][0]

		type d struct {
			Text      string "xml:\",chardata\""
			IDDomanda string "xml:\"idDomanda,attr\""
			Tipo      string "xml:\"tipo,attr\""
			Opzioni   struct {
				Text    string   "xml:\",chardata\""
				Opzione []string "xml:\"opzione\""
			} "xml:\"opzioni\""
		}

		for k, v := range m {

			if _, err := strconv.Atoi(k); err == nil {

				var t d
				t.IDDomanda = k

				t.Text = v[0]

				t.Tipo = v[2]
				if t.Tipo == "multipla" {
					t.Opzioni.Opzione = v[3:]
				}
				switch v[1] {
				case "adoption":
					if !stringInSlice(newsurvey.Department, m["Funzioni interessate"]) {
						continue
					}
					newsurvey.Domande.Adoption = append(newsurvey.Domande.Adoption, t)
				default:
					newsurvey.Domande.Domanda = append(newsurvey.Domande.Domanda, t)
				}
			}
		}

		inizio, _ := time.Parse("20060102", newsurvey.Inizio)
		fine, _ := time.Parse("20060102", newsurvey.Fine)
		newsurvey.Inizio = inizio.Format("2006-01-02")
		newsurvey.Fine = fine.Format("2006-01-02")

		// Serve template
		err := templates.ExecuteTemplate(w, "survey.gohtml", newsurvey)
		if err != nil {
			log.Println(err)
		}

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Println(r.Form)

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

func readSheet(sheetID string) []byte {
	os.Setenv("HTTPS_PROXY", httpsproxy)

	urlFunction := "https://europe-west6-ctio-8274d.cloudfunctions.net/SheetRead?sheetID=" + sheetID + "&readRange=A1:AA"
	req, err := http.NewRequest("GET", urlFunction, nil)

	proxyURL, err := url.Parse(httpsproxy)
	// myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Errore nel recupero dati da sheet %s : %v", sheetID, err)
	}

	return body

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

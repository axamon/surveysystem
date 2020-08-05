package main

import (
	"log"
	"net/http"
	"strings"
)

// index serve il template iniziale dell'applicazione.
func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		ua := r.UserAgent()
		if !strings.Contains(ua, "Chrome") && !strings.Contains(ua, "Firefox") {
			log.Printf("Browser non supportato: %s\n", ua)
			err := templates.ExecuteTemplate(w, "nobrowser.gohtml", nil)
			if err != nil {
				log.Println(err)
			}
			return
		}

		err := templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}

}

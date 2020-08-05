package main

import (
	"fmt"
	"log"
	"net/http"
)

// index serve il template iniziale dell'applicazione.
func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		ua := r.UserAgent()

		fmt.Println(ua)

		err := templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}

}

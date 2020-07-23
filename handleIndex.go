package main

import (
	"log"
	"net/http"
)

// index serve il template iniziale dell'applicazione.
func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		// err :=templates["index"].Execute(w, nil)
		err := templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

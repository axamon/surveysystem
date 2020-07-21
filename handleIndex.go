package main

import (
	"html/template"
	"log"
	"net/http"
)

// index serve il template iniziale dell'applicazione.
func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		var indexTmpl = template.Must(template.ParseFiles(
			"templates/index.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml"))
		// err :=templates["index"].Execute(w, nil)
		err := indexTmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

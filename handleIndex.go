package main

import (
	"log"
	"net/http"
)

// index serve il template iniziale dell'applicazione.
func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
<<<<<<< HEAD

		//var indexTmpl =
		// err :=templates["index"].Execute(w, nil)
=======
>>>>>>> b2ff8f323abf46180ced9a5358a03b91dfc82293
		err := templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}

}

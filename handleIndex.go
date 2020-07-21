package main

import (
	"html/template"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var indexTmpl = template.Must(template.ParseFiles("templates/index.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))
	// err :=templates["index"].Execute(w, nil)
	err := indexTmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

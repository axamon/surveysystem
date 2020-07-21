package main

import (
	"html/template"
	"log"
	"net/http"
)

var errTmpl = template.Must(template.ParseFiles("templates/error.gohtml", "templates/header.gohtml", "templates/footer.gohtml"))

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.RequestURI() {
		case "/survey":
			session, err := store.Get(r, "surveyCTIO")
			if err != nil {
				log.Println(err)
			}
			// Se l'utente non è autenticato restituisce il template errore.
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				w.WriteHeader(http.StatusForbidden)
				err := errTmpl.Execute(w, nil)
				if err != nil {
					log.Println(err)
				}
				return
			}

			fallthrough
		default:
			next.ServeHTTP(w, r)
			return
		}
	})
}

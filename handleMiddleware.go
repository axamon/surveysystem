package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		switch r.URL.RequestURI() {
		case "/survey":
			session, err := store.Get(r, "surveyCTIO")
			if err != nil {
				log.Println(err)
			}
			// Se l'utente non è autenticato restituisce il template errore.
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				w.WriteHeader(http.StatusForbidden)
				err := templates.ExecuteTemplate(w, "error.gohtml", nil)
				if err != nil {
					log.Println(err)
				}
				return
			}
			select {
			case <-ctx.Done():
				log.Printf("Took too long to serve %s: %v\n", r.RequestURI, ctx.Err())
			}

			fallthrough
		default:
			next.ServeHTTP(w, r)
			return
		}
	})
}

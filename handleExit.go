package main

import (
	"log"
	"net/http"
	"time"
)

func exit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		session, err := store.Get(r, "surveyCTIO")
		if err != nil {
			log.Println(err)
		} else {
			// Revoke users authentication
			session.Values["authenticated"] = false
			session.Options.MaxAge = -1
			err = session.Save(r, w)
			if err != nil {
				log.Println(err)
			}
			if _, ok := session.Values["utente"]; ok {
				log.Printf("Logout effettuato per %s", session.Values["utente"].(string))
			}

		}

		err = templates.ExecuteTemplate(w, "exit.gohtml", footerData)
		if err != nil {
			log.Println(err)
		}
		go func() {
			time.Sleep(1 * time.Second)
			log.Println("Uscita da app.")
		}()

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

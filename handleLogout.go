package main

import (
	"log"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
<<<<<<< HEAD
	case "GET":
=======
>>>>>>> b2ff8f323abf46180ced9a5358a03b91dfc82293

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

		err = templates.ExecuteTemplate(w, "logout.gohtml", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

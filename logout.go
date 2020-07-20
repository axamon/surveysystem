package main

import (
	"log"
	"net/http"
	"text/template"
)

func logout(w http.ResponseWriter, r *http.Request) {
	var logoutTmpl = template.Must(template.ParseFiles("templates/logout2.gohtml", "templates/footer.gohtml"))

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

	err = logoutTmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

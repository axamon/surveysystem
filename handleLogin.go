package main

import (
	"fmt"
	"log"
	"net/http"
	"survey/ldaplogin"
)

// login serve a gestire e verificare l'autenicazione e
// autorizzazione utente.
func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "surveyCTIO")
	if err != nil {
		log.Println(err)
	}

	switch r.Method {

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)

		matricola := r.FormValue("username")
		password := r.FormValue("password")

		ok, nomeCognome, err := ldaplogin.IsOK(matricola, password)
		if err != nil {
			//http.Redirect(w, r, "/login", 301)
			log.Println(err)
		}
		// ripulisce la passoword per non farla girare
		password = "******"
		// Set user as authenticated
		if ok {
			session.Values["authenticated"] = true
			session.Values["matricola"] = matricola
			session.Values["utente"] = nomeCognome
		}
		// permette ad "Admin" di entrare senza verifica LDAP.
		if matricola == "Admin" {
			session.Values["authenticated"] = true
			session.Values["matricola"] = "Admin"
			session.Values["utente"] = "Admin"
		}
		session.Save(r, w)
		r.Method = "GET"
		r.RequestURI = "/survey"
		survey(w, r)
	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

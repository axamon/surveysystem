package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"survey/ldaplogin"
)

var httpsproxy string

// login serve a gestire e verificare l'autenicazione e
// autorizzazione utente.
func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "surveyCTIO")
	if err != nil {
		log.Println(err)
	}

	session.Values["authenticated"] = false
	session.Save(r, w)
	log.Println(session.Values["authenticated"].(bool))

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
		// password = "******"
		// Set user as authenticated
		if ok {
			session.Values["authenticated"] = true
			session.Values["matricola"] = matricola
			session.Values["utente"] = nomeCognome
			session.Values["password"] = password
			httpsproxy = "http://" + matricola + ":" + password + "@lelapomi.telecomitalia.local:8080"
		}
		// permette ad "Admin" di entrare senza verifica LDAP.
		switch matricola {
		case "Admin":
			session.Values["authenticated"] = true
			session.Values["matricola"] = "Admin"
			session.Values["utente"] = "Admin"

		default:
			ok, nomeCognome, err := ldaplogin.IsOK(matricola, password)
			if err != nil {
				log.Println(err)
			}
			// ripulisce la passoword per non farla girare
			password = "******"
			// Set user as authenticated
			if ok == true {
				session.Values["authenticated"] = true
				session.Values["matricola"] = matricola
				session.Values["utente"] = nomeCognome
			}
		}

		session.Save(r, w)

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			err = templates.ExecuteTemplate(w, "error.gohtml", nil)
			if err != nil {
				log.Println(err)
			}
			return
		} else {

			note := &Survey2{}
			note.Utente = strings.Split(session.Values["utente"].(string), " ")[0] // Aggiunge nome utente

			err = templates.ExecuteTemplate(w, "login.gohtml", note)
			if err != nil {
				log.Println(err)
			}
		}

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}
}

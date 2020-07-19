package main

import (
	"log"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "surveyCTIO")
	if err != nil {
		log.Println(err)
	}
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Logout effettuato per %s", session.Values["utente"].(string))
	// w.Write([]byte("Logout effettuato\n"))
	// time.Sleep(2 * time.Second)
	r.Method = "GET"
	http.Redirect(w, r, "/", http.StatusFound)
}

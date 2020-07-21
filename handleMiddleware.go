package main

import (
	"net/http"
	"strings"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "static") {
			return
		}
		session, _ := store.Get(r, "surveyCTIO")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "vietato", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package main

import "net/http"

func main() {

	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/", fs)

	http.ListenAndServe(":8000", r)
}

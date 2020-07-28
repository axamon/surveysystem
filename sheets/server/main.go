package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	r := http.NewServeMux()

	r.HandleFunc("/", index)

	http.ListenAndServe(":8080", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	var p Answers

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
	}

	sheetID := p.SheetID
	val := p.Val

	log.Println(sheetID, val)

}

// Answers sono le risposte da inviare.
type Answers struct {
	SheetID string `json:"sheetID"`
	Foglio  string `json:"foglio"`
	Val     string `json:"val"`
}

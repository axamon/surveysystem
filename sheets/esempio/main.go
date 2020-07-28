package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	// var sheetID = flag.String("sheetID", "1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo", "ID del foglio google da usare")
	// var foglio = flag.String("foglio", "Risposte", "Il foglio da usare")
	// flag.Parse()

	data := []string{"uno", "due", "tre"}

	writeToCSV(data)

}

// writeToCSV registra le risposte nel file csv.
func writeToCSV(data []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		err           error
		list, sheetID string
	)

	list = strings.Join(data, ";")

	sheetID = "1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo"

	answers := new(Answers)
	answers.SheetID = sheetID
	answers.Val = list

	payload, err := json.Marshal(answers)
	if err != nil {
		log.Printf("marshal of answers in error: %v\n", err)
	}

	url := "https://us-central1-ctio-8274d.cloudfunctions.net/SheetAppend"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println(resp.Status, err, resp.Body)

	// link visualizzazione risultati
	//https://docs.google.com/spreadsheets/d/1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo/edit?usp=sharing

	select {
	case <-ctx.Done():
		log.Printf("Save to sheet took too long: %v\n", ctx.Err())
	}
	return err
}

// Answers sono le risposte da inviare.
type Answers struct {
	SheetID string `json:"sheetID"`
	Val     string `json:"val"`
}

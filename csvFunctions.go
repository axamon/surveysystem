package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// writeToCSV registra le risposte nel file csv.
func writeToCSV(data map[string][]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		err                      error
		record                   []string
		matricola, list, sheetID string
	)

	for i := 1; i < len(data)-2; i++ {
		if v, ok := data[strconv.Itoa(i)]; ok {
			record = append(record, strings.Join(v, ","))
		} else {
			record = append(record, "")
		}
	}
	matricola = "\"" + strings.Join(data["matricola"], "") + "\""
	record = append(record,
		time.Now().Format("20060102T15:04"),
		matricola)

	list = strings.Join(record, ";")

	sheetID = "1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo"

	answers := new(Answers)
	answers.SheetID = sheetID
	answers.Val = list

	payload, err := json.Marshal(answers)
	if err != nil {
		log.Printf("marshal of answers in error: %v\n", err)
	}

	url := "https://europe-west6-ctio-8274d.cloudfunctions.net/SheetAppend"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("X-Custom-Header", "answers")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// link visualizzazione risultati
	// https://docs.google.com/spreadsheets/d/1KXUdTBXDhGvBU1U8SKuf1OBUqYpyQdLW6GMHTxylk2Y/edit#gid=0
	os.Setenv("HTTPS_PROXY", httpsproxy)

	request := gorequest.New()
	_, _, errs := request.Proxy(httpsproxy).Get("https://us-central1-ctio-8274d.cloudfunctions.net/SheetAppend?val=" + encoded + "&sheetID=" + sheetID + "&foglio=Risposte").End()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	select {
	case <-ctx.Done():
		log.Printf("Save to sheet took too long: %v\n", ctx.Err())
	}
	return err
}

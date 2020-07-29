package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// writeToCSV registra le risposte nel file csv.
func writeToCSV(data map[string][]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		err       error
		record    []string
		matricola string
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

	answers := new(Answers)
	answers.SheetID = "1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo"
	answers.Val = strings.Join(record, ";")

	payload, err := json.Marshal(answers)
	if err != nil {
		log.Printf("marshal of answers in error: %v\n", err)
	}

	os.Setenv("HTTPS_PROXY", httpsproxy)

	urlFunction := "https://europe-west6-ctio-8274d.cloudfunctions.net/SheetAppend"
	req, err := http.NewRequest("POST", urlFunction, bytes.NewBuffer(payload))
	req.Header.Set("X-Custom-Header", "answers")
	req.Header.Set("Content-Type", "application/json")

	proxyURL, err := url.Parse(httpsproxy)
	// myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		log.Printf("Save to sheet took too long: %v\n", ctx.Err())
	}
	return err
}

package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// writeToCSV registra le risposte nel file csv.
func writeToCSV(data map[string][]string) error {

	var (
		err                               error
		record                            []string
		matricola, list, encoded, sheetID string
	)

	for i := 1; i < len(data)-2; i++ {
		if v, ok := data[strconv.Itoa(i)]; ok {
			record = append(record, strings.Join(v, ","))
		} else {
			record = append(record, "")
		}
	}
	matricola = "\"" + strings.Join(data["matricola"], "") + "\""
	record = append(record, time.Now().Format("20060102T15:04"))
	record = append(record, matricola)

	list = strings.Join(record, ";")
	encoded = base64.StdEncoding.EncodeToString([]byte(list))

	sheetID = "1KXUdTBXDhGvBU1U8SKuf1OBUqYpyQdLW6GMHTxylk2Y"
	// link visualizzazione risultati
	// https://docs.google.com/spreadsheets/d/1KXUdTBXDhGvBU1U8SKuf1OBUqYpyQdLW6GMHTxylk2Y/edit#gid=0
	_, err = http.Get("https://us-central1-ctio-8274d.cloudfunctions.net/SheetAppend?val=" + encoded + "&sheetID=" + sheetID)

	return err
}

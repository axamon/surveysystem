package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var csvlock sync.RWMutex

// writeToCSV registra le risposte nel file csv.
func writeToCSV(data map[string][]string) error {
	csvlock.Lock()
	defer csvlock.Unlock()

	fileCSV := "surveyID" + strings.Join(data["surveyID"], "") + ".csv"
	f, err := os.OpenFile(fileCSV, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	defer csvwriter.Flush()

	var record []string
	for i := 1; i < len(data)-2; i++ {
		if v, ok := data[strconv.Itoa(i)]; ok {
			record = append(record, strings.Join(v, ","))
		} else {
			record = append(record, "")
		}
	}
	matricola := "\"" + strings.Join(data["matricola"], "") + "\""
	record = append(record, time.Now().Format("20060102T15:04"))
	record = append(record, matricola)
	fmt.Println(len(record))
	for i, r := range record {
		fmt.Println(i, r)
	}

	for k, v := range data {
		fmt.Println(k, v)
	}
	csvwriter.Comma = ';'

	err = csvwriter.Write(record)
	if err != nil {
		return err
	}

	return nil
}

// createFileCsv crea il file Csv che contiene le risposte.
func createFileCsv(path string, fields int) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		file.Close()
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		defer f.Close()

		var csvheader []string
		for i := 0; i < fields; i++ {
			if i == 0 {
				csvheader = append(csvheader, "#Domanda"+strconv.Itoa(i+1))
				continue
			}
			csvheader = append(csvheader, "Domanda"+strconv.Itoa(i+1))
		}
		csvheader = append(csvheader, "timestamp", "matricola")
		w := csv.NewWriter(f)
		w.Comma = ';'
		w.Write(csvheader)
		w.Flush()

		fmt.Println("Created file: ", path)
	}

	return nil
}

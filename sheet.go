package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func sheet(sheetID string) {

	sheetID = "1KXUdTBXDhGvBU1U8SKuf1OBUqYpyQdLW6GMHTxylk2Y"

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v\n", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

	writeRange := "A1"

	var vr sheets.ValueRange

	myval := []interface{}{"One", "Two", "Three"}
	vr.Values = append(vr.Values, myval)

	_, err = srv.Spreadsheets.Values.Append(sheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println(err)
	}

}

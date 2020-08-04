package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// SheetRead reads data from a google sheet.
func SheetRead(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		values := r.URL.Query()

		sheetID := values.Get("sheetID")
		readRange := values.Get("readRange")

		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			fmt.Fprintf(w, "Unable to read client secret file: %v", err)
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			fmt.Fprintf(w, "Unable to parse client secret file to config: %v", err)
		}
		client := getClient(config)

		srv, err := sheets.New(client)
		if err != nil {
			fmt.Fprintf(w, "Unable to retrieve Sheets client: %v", err)
		}

		resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
		if err != nil {
			fmt.Fprintf(w, "Unable to retrieve data from sheet: %v", err)
		}

		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			for _, row := range resp.Values {
				// Print columns A and E, which correspond to indices 0 and 4.
				for _, element := range row {
					fmt.Fprintf(w, "%s,", element)
				}
			}
		}

	default:
		http.Error(w, "metodo non permesso", http.StatusMethodNotAllowed)
	}
}

// SheetAppend appends data to a google sheet.
func SheetAppend(w http.ResponseWriter, r *http.Request) { //sheetID string, val []string) {

	var a Answers

	switch r.Method {
	case "POST":

		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		sheetID := a.SheetID
		val := a.Val

		fmt.Fprint(w, sheetID, val)

	default:
		http.Error(w, "Metodo non permesso", http.StatusMethodNotAllowed)
	}

	err := appendToSheet(a)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}

func appendToSheet(a Answers) error {

	var srv *sheets.Service
	var client *http.Client
	var config *oauth2.Config
	var list []string
	var writeRange string
	var myval []interface{}
	var vr sheets.ValueRange

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		err = fmt.Errorf("Unable to read client secret file: %v", err)
		goto ERR
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err = google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		err = fmt.Errorf("Unable to parse client secret file to config: %v", err)
		goto ERR
	}
	client = getClient(config)

	srv, err = sheets.New(client)
	if err != nil {
		err = fmt.Errorf("Unable to retrieve Sheets client: %v", err)
		goto ERR
	}

	writeRange = "A1"

	if a.Foglio != "" {
		writeRange = a.Foglio + "!" + writeRange
	}

	list = strings.Split(a.Val, ";")

	myval = []interface{}{}
	for _, l := range list {
		myval = append(myval, l)
	}
	vr.Values = append(vr.Values, myval)
	vr.Values = append(vr.Values, nil)

	_, err = srv.Spreadsheets.Values.Append(a.SheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		err = fmt.Errorf("Impossibile effettuare append su google sheet %s : %v", a.SheetID, err)
		goto ERR
	}

ERR:
	return err
}

// SheetAppendLocal appends data to a google sheet.
func SheetAppendLocal(sheetID, foglio string, val []string) {

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

	//writeRange = foglio + "!" + writeRange

	var vr sheets.ValueRange

	list := val

	myval := []interface{}{}
	for _, l := range list {
		myval = append(myval, l)
	}
	vr.Values = append(vr.Values, myval)
	vr.Values = append(vr.Values, nil)

	_, err = srv.Spreadsheets.Values.Append(sheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println(err)
	}

}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Answers sono le risposte degli utenti.
type Answers struct {
	SheetID string `json:"sheetID"`
	Foglio  string `json:"foglio"`
	Val     string `json:"val"`
}

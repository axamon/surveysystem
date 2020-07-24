package main

import "encoding/xml"

// Survey è la struttura in cui parsare le domande
type Survey struct {
	Utente    string
	Matricola string
	XMLName   xml.Name `xml:"survey"`
	Text      string   `xml:",chardata"`
	ID        string   `xml:"id,attr"`
	Titolo    string   `xml:"titolo,attr"`
	Inizio    string   `xml:"inizio,attr"`
	Fine      string   `xml:"fine,attr"`
	Domande   struct {
		Text    string `xml:",chardata"`
		Domanda []struct {
			Text      string `xml:",chardata"`
			IDDomanda string `xml:"idDomanda,attr"`
			Tipo      string `xml:"tipo,attr"`
			Opzioni   struct {
				Text    string   `xml:",chardata"`
				Opzione []string `xml:"opzione"`
			} `xml:"opzioni"`
		} `xml:"domanda"`
	} `xml:"domande"`
}

type Survey2 struct {
	Utente    string
	Matricola string
	XMLName   xml.Name `xml:"survey"`
	Text      string   `xml:",chardata"`
	ID        string   `xml:"id,attr"`
	Titolo    string   `xml:"titolo,attr"`
	Inizio    string   `xml:"inizio,attr"`
	Fine      string   `xml:"fine,attr"`
	Video     string   `xml:"video,attr"`
	Domande   struct {
		Text    string `xml:",chardata"`
		Domanda []struct {
			Text      string `xml:",chardata"`
			IDDomanda string `xml:"idDomanda,attr"`
			Tipo      string `xml:"tipo,attr"`
			Opzioni   struct {
				Text    string   `xml:",chardata"`
				Opzione []string `xml:"opzione"`
			} `xml:"opzioni"`
		} `xml:"domanda"`
	} `xml:"domande"`
}

// Utente contiene le informazioni anagrafiche.
type Utente struct {
	Matricola string
	Nome      string
	Cognome   string
	Mail      string
	Surveys   []Survey
	Risposte  [][]string
}

// InternalTemplate are templates made from embedded data.
type InternalTemplate struct {
	Name string
	Text string
}

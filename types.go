package main

import "encoding/xml"

// Survey3 è la struttura in cui parsare le domande
type Survey3 struct {
	TimestampInizio string
	Utente          string
	Matricola       string
	Department      string
	Versione        string
	XMLName         xml.Name `xml:"survey"`
	Text            string   `xml:",chardata"`
	ID              string   `xml:"id,attr"`
	Titolo          string   `xml:"titolo,attr"`
	Inizio          string   `xml:"inizio,attr"`
	Fine            string   `xml:"fine,attr"`
	Video           string   `xml:"video,attr"`
	Domande         struct {
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
		Adoption []struct {
			Text      string `xml:",chardata"`
			IDDomanda string `xml:"idDomanda,attr"`
			Tipo      string `xml:"tipo,attr"`
			Opzioni   struct {
				Text    string   `xml:",chardata"`
				Opzione []string `xml:"opzione"`
			} `xml:"opzioni"`
		} `xml:"adoption"`
	} `xml:"domande"`
}

// Answers sono le risposte degli utenti.
type Answers struct {
	SheetID string `json:"sheetID"`
	Foglio  string `json:"foglio"`
	Val     string `json:"val"`
}

// FooterInfo informazioni da mostrare nel footer
type FooterInfo struct {
	Anno     string
	Autore   string
	Versione string
}

// Surveys lista dei survey esistenti.
type Surveys struct {
	NomeUtente   string
	ListaSurveys []Survey3
	Versione     string
}

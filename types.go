package main

import "encoding/xml"

// Survey è la struttura in cui parsare le domande
type Survey struct {
	XMLName xml.Name `xml:"survey"`
	Text    string   `xml:",chardata"`
	Titolo  string   `xml:"titolo,attr"`
	Inizio  string   `xml:"inizio,attr"`
	Fine    string   `xml:"fine,attr"`
	Domande struct {
		Text    string `xml:",chardata"`
		Domanda []struct {
			Text    string `xml:",chardata"`
			Tipo    string `xml:"tipo,attr"`
			Riposte struct {
				Text    string   `xml:",chardata"`
				Riposta []string `xml:"riposta"`
			} `xml:"riposte"`
		} `xml:"domanda"`
	} `xml:"domande"`
}

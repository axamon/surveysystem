package main

import (
	"bufio"
	"fmt"
	"strings"
)

var text = `Nome Survey,Cyber Security Process Adoption,
SurveyID,1,
Video,https://web.microsoftstream.com/embed/video/4c88d6fb-f9ad-4ca1-9a4e-4d51cf77e5f7?autoplay=false&amp;showinfo=true,
Inizio,20200731,
Fine,20200815,
DomandaID,Quesito,tipo,opzioni,
Domanda1,Quanti sono i capisaldi del sistema Crisis Management di TIM?,libera,
Domanda2,Aumentare la resilienza è uno degli obiettivi della procedura?,booleana,
Domanda3,La pianificazione rientra nelle 4 fasi del Crisis Management di TIM?,booleana,
Domanda4,Da quante funzioni è composto l'Operationa Crisis Team in CTIO?,libera,
Domanda5,Che livello di gravità viene assegnato a un evento critico?,libera,
Domanda6,L'OCT deve occuparsi anche delle relazioni con i social media?,booleana,
Domanda7,Massimo quante ore sono previste da KPI per l'attivazione operativa?,libera,
Domanda8,Quali sono gli attori che partecipano al processo?,multipla,Operationl Crisis Team,Crisis Chief,CEO,Security Operation Center,Tutte le strutture di 2° livello,Crisis Management,Crisis Managemen Committee,`

func main() {

	s := bufio.NewScanner(strings.NewReader(text))
	// var tsv = make(map[string][]string)

	var m = make(map[string][]string)

	for s.Scan() {
		// fmt.Println(s.Text())

		list := strings.Split(s.Text(), ",")

		m[list[0]] = list[1:]

	}

	fmt.Println(m["Nome Survey"][0])

}

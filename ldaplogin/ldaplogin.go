// Package ldaplogin enables to check if users and password
// match on a LDAP server database.
package ldaplogin

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

// LDAPURL è la URL del server LDAP aziendale.
var LDAPURL = "ldap://directory.cww.telecomitalia.it:389"

// LDAPBASE è lo starting point per la ricerca.
var LDAPBASE = "OU=Telecomitalia,O=Telecom Italia Group"

// IsOK verifica se le credenziali passate sono accettate.
func IsOK(username, password string) (bool, UserInfo, error) {

	var (
		userdn, usercn, department string
		err                        error
		searchRequest              *ldap.SearchRequest
		ldapResponse               *ldap.SearchResult
	)

	// searchRequest è la richista da inviare al server LDAP.
	searchRequest = ldap.NewSearchRequest(
		LDAPBASE,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn", "department"}, // viene richiesta matricola con nome e cognome.
		nil,
	)

	// Si collega al server LDAP.
	ldapconn, err := ldap.DialURL(LDAPURL)
	if err != nil {
		err = fmt.Errorf("Impossibile connettersi a LDAP: %v", err)
		goto ERR
	}

	// ldapResponse è la risposta del server LDAP.
	ldapResponse, err = ldapconn.Search(searchRequest)
	if err != nil {
		goto ERR
	}

	// Se non si ottiene una sola risposta esce con errore.
	if len(ldapResponse.Entries) != 1 {
		err = fmt.Errorf("utente %s non trovato", username)
		goto ERR
	}

	// userdn è il nome interno a LDAP che ideticfica l'utente.
	userdn = ldapResponse.Entries[0].DN
	usercn = ldapResponse.Entries[0].GetAttributeValue("cn")
	department = ldapResponse.Entries[0].GetAttributeValue("department")

	// tenta il binding su LDAP con username valido su LDAP e password.
	err = ldapconn.Bind(userdn, password)
	if err != nil {
		err = fmt.Errorf("impossibile loggarsi: %v", err)
		goto ERR
	}

ERR:
	return err == nil, UserInfo{Matricola: userdn, NomeCognome: usercn, Department: department}, err
}

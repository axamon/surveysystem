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
func IsOK(username, password string) (bool, string, error) {

	ldapconn, err := ldap.DialURL(LDAPURL)
	if err != nil {
		// ! Se non riesce logga l'errore.
		return false, "", fmt.Errorf("Impossibile connettersi a LDAP: %v", err)
	}

	// searchRequest è la richista da inviare al server LDAP.
	searchRequest := ldap.NewSearchRequest(
		LDAPBASE,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn"}, // viene richiesto solo il nome.
		nil,
	)

	// ldapResponse è la risposta del server LDAP.
	ldapResponse, err := ldapconn.Search(searchRequest)
	if err != nil {
		return false, "", err
	}

	// Se non si ottiene una sola risposta esce con errore.
	if len(ldapResponse.Entries) != 1 {
		return false, "", fmt.Errorf("utente %s non trovato", username)
	}

	// userdn è il nome interno a LDAP che ideticfica l'utente.
	userdn := ldapResponse.Entries[0].DN
	usercn := ldapResponse.Entries[0].GetAttributeValue("cn")

	// tenta il binding su LDAP con username valido su LDAP e password.
	err = ldapconn.Bind(userdn, password)
	if err != nil {
		return false, "", fmt.Errorf("impossibile loggarsi: %v", err)
	}

	return true, usercn, nil
}

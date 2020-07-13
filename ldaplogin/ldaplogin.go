package ldaplogin

import (
        "fmt"
        "log"

        "github.com/go-ldap/ldap"
)

// LDAPURL è la URL del server LDAP aziendale.
var LDAPURL = "ldap://directory.cww.telecomitalia.it:389"

// LDAPBASE è lo starting point per la ricerca.
var LDAPBASE = "OU=Telecomitalia,O=Telecom Italia Group"

// ldapconn è la sessione da utilizzare per la ricerca su LDAP.
var ldapconn *ldap.Conn

func init() {
        // Apre una connessione con LDAP aziendale.
        l, err := ldap.DialURL(LDAPURL)
        if err != nil {
                // ! Se non riesce esce con errore.
                log.Fatal("Impossibile connettersi a LDAP: ", err)
        }
        ldapconn = l
}

// IsOK verifica se le credenziali passate sono accettate.
func IsOK(username, password string) (bool, error) {

        // searchRequest è la richista da inviare al server LDAP.
        searchRequest := ldap.NewSearchRequest(
                LDAPBASE,
                ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
                fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
                []string{"dn"}, // viene richiesto solo il nome.
                nil,
        )

        // ldapResponse è la risposta del server LDAP.
        ldapResponse, err := ldapconn.Search(searchRequest)
        if err != nil {
                return false, err
        }

        // Se non si ottiene una sola risposta esce con errore.
        if len(ldapResponse.Entries) != 1 {
                return false, fmt.Errorf("utente %s non trovato", username)
        }

        // userdn è il nome interno a LDAP che ideticfica l'utente.
        userdn := ldapResponse.Entries[0].DN

        // tenta il binding su LDAP con username valido su LDAP e password.
        err = ldapconn.Bind(userdn, password)
        if err != nil {
                return false, fmt.Errorf("impossibile loggarsi: %v", err)
        }

        return true, nil
}

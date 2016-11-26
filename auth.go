package main

import (
	"encoding/json"
	"log"
)

const (
	authFileName = "auth.db"
)

type AuthToken struct {
	Token string `json:token`
	Email string `json:email`
}

type TokenDB []AuthToken

func makeTokenDB() TokenDB {
	var tokens TokenDB
	authDB := tryReadFile(authFileName)
	err := json.Unmarshal(authDB, &tokens)
	if err != nil {
		log.Println("error reading auth token db:", err)
	}
	log.Printf("found %d auth tokens\n", len(tokens))
	log.Printf("%v\n", tokens)
	return tokens
}

func (db TokenDB) findToken(token string) (email string) {
	for _, i := range db {
		if i.Token == token {
			email = i.Email
			return
		}
	}
	return
}

// isAuthorized tries to find the auth token given in entry.
// It will the change the entry parameter by replacing the auth
// token with the associated email address. This is to have the
// auth token not end up in the secret database.
func (db TokenDB) isAuthorized(entry *StoreEntry) bool {
	email := db.findToken(entry.AuthToken)
	if email == "" {
		return false
	}
	entry.AuthToken = email
	return true
}

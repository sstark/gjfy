package main

import (
	"encoding/json"
)

const authDB = `[{"token":"test","email":"test@example.org"}]`

type AuthToken struct {
	Token string `json:token`
	Email string `json:email`
}

type TokenDB []AuthToken

func makeTokenDB() TokenDB {
	tokens := make(TokenDB, 0)
	json.Unmarshal([]byte(authDB), &tokens)
	return tokens
}

func (db *TokenDB) findToken(token string) (email string) {
	for _, i := range *db {
		if i.Token == token {
			email = i.Email
			return
		}
	}
	return
}

func isAuthorized(entry *StoreEntry) bool {
	tokens := makeTokenDB()
	email := tokens.findToken(entry.AuthToken)
	if email == "" {
		return false
	}
	entry.AuthToken = email
	return true
}

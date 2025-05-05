package tokendb

import (
	"encoding/json"
	"log"

	"github.com/sstark/gjfy/store"
)

const (
	AuthFileName = "auth.db"
)

type AuthToken struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type TokenDB []AuthToken

func MakeTokenDB(b []byte) TokenDB {
	var tokens TokenDB
	err := json.Unmarshal(b, &tokens)
	if err != nil {
		log.Println("error reading auth token db:", err)
	}
	for i, entry := range tokens {
		if entry.Token == "" {
			log.Printf("token field empty or missing in entry #%d", i)
			return nil
		}
		if entry.Email == "" {
			log.Printf("email field empty or missing in entry #%d", i)
			return nil
		}
	}
	log.Printf("found %d auth tokens\n", len(tokens))
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

// IsAuthorized tries to find the auth token given in entry.
// It will then change the entry parameter by replacing the auth
// token with the associated email address. This is to have the
// auth token not end up in the secret database.
func (db TokenDB) IsAuthorized(entry *store.StoreEntry) bool {
	email := db.findToken(entry.AuthToken)
	if email == "" {
		return false
	}
	entry.AuthToken = email
	return true
}

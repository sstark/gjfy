package httpio

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/sstark/gjfy/store"
	"github.com/sstark/gjfy/tokendb"
)

const (
	maxData = 1048576 // 1MB
)

type jsonError struct {
	Error string `json:"error"`
}

func HandleApiGet(memstore store.SecretStore, urlbase string, fNotify bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len(ApiGet):]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if entry, ok := memstore.GetEntryInfo(id, urlbase, Get, ApiGet); !ok {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("entry not found: %s", id)
			if jerr := json.NewEncoder(w).Encode(jsonError{"not found"}); jerr != nil {
				panic(jerr)
			}
		} else {
			memstore.Click(id, r, fNotify)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(entry); err != nil {
				panic(err)
			}
		}
	})
}

func HandleApiNew(memstore store.SecretStore, urlbase string, auth *tokendb.TokenDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var entry store.StoreEntry

		body, err := io.ReadAll(io.LimitReader(r.Body, maxData))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.Unmarshal(body, &entry); err != nil {
			w.WriteHeader(422) // unprocessable entity
			log.Printf("error processing json: %s", err)
			if jerr := json.NewEncoder(w).Encode(jsonError{err.Error()}); jerr != nil {
				panic(jerr)
			}
		} else if !auth.IsAuthorized(&entry) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("unauthorized try to make new entry")
			if jerr := json.NewEncoder(w).Encode(jsonError{"unauthorized"}); jerr != nil {
				panic(jerr)
			}
		} else {
			id := memstore.AddEntry(entry, "")
			newEntry, _ := memstore.GetEntryInfoHidden(id, urlbase, Get, ApiGet)
			log.Println("New ID:", id)
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(newEntry); err != nil {
				panic(err)
			}
		}
	})
}

package httpio

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"

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
		id := path.Base(r.URL.Path)
		if entry, ok := memstore.GetEntryInfo(id, urlbase, Get, ApiGet); !ok {
			log.Printf("entry not found: %s", id)
			jsonRespond(w, http.StatusNotFound, jsonError{"not found"})
		} else {
			memstore.Click(id, r, fNotify)
			jsonRespond(w, http.StatusOK, entry)
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
		if err := json.Unmarshal(body, &entry); err != nil {
			log.Printf("error processing json: %s", err)
			jsonRespond(w, http.StatusUnprocessableEntity, jsonError{err.Error()})
		} else if !auth.IsAuthorized(&entry) {
			log.Printf("unauthorized when trying to make new entry")
			jsonRespond(w, http.StatusUnauthorized, jsonError{"unauthorized"})
		} else {
			id := memstore.AddEntry(entry, "")
			newEntry, _ := memstore.GetEntryInfoHidden(id, urlbase, Get, ApiGet)
			log.Println("New ID:", id)
			jsonRespond(w, http.StatusCreated, newEntry)
		}
	})
}

func jsonRespond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	// FIXME: don't panic
	if jerr := json.NewEncoder(w).Encode(data); jerr != nil {
		panic(jerr)
	}
}

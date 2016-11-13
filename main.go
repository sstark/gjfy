package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	//"html"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	schemeHost = "http://localhost"
	listen     = ":9154"
	uApiGet    = "/api/v1/get/"
	uApiNew    = "/api/v1/new"
	uGet       = "/g"
	maxData    = 1048576 // 1MB
)

// In-memory representation of a secret.
type StoreEntry struct {
	Secret    string    `json:"secret"`
	MaxClicks int       `json:"max_clicks"`
	DateAdded time.Time `json:"date_added"`
}

// Secret augmented with computed fields.
type StoreEntryInfo struct {
	StoreEntry
	Id        string `json:"id"`
	PathQuery string `json:"path_query"`
	Url       string `json:"url"`
	Clicks    int    `json:"clicks"`
}

type secretStore map[string]StoreEntry

// hashStruct returns a hash from an arbitrary structure, usable in a URL.
func hashStruct(data interface{}) (hash string) {
	hashBytes := sha256.Sum256([]byte(fmt.Sprintf("%#v", data)))
	hash = base64.URLEncoding.EncodeToString(hashBytes[:])
	return
}

// NewEntry adds a new secret to the store.
func (st secretStore) NewEntry(e StoreEntry) string {
	e.DateAdded = time.Now()
	id := hashStruct(e)
	st[id] = e
	return id
}

// GetEntry retrives a secret from the store.
func (st secretStore) GetEntry(id string) (se StoreEntry, ok bool) {
	se, ok = st[id]
	return
}

// GetEntryInfo wraps GetEntry and adds some computed fields.
func (st secretStore) GetEntryInfo(id string) (si StoreEntryInfo, ok bool) {
	entry, ok := st.GetEntry(id)
	pathQuery := uGet + "?" + id
	url := schemeHost + listen + pathQuery
	return StoreEntryInfo{entry, id, pathQuery, url, 3}, ok
}

func main() {
	store := make(secretStore)
	store["1234"] = StoreEntry{"geheim", 5, time.Now()}

	/*
		/							# intro page
		/api/v1
			/api/v1/get/34g34g243	# get entry, decrease counter
			/api/v1/new				# generate new entry (takes POST data)
			/api/v1/info/34g34g243	# show info for entry
		/g?34g34g243				# show entry, decrease counter
		/n							# generate new entry (takes POST data)
		/i?34g34g243				# show info for entry (e. g. after creating)

		"entryinfo": {
			"id": "34g34g243",
			"path_query": "/g?34g34g243",
			"url": "https://klause.bla.de/g?34g34g243",
			"date_added": <jsondate>,
			"valid_until": <jsondate>,
			"clicks": 0,
			"max_clicks": 1,
		}

	*/

	http.HandleFunc(uApiGet, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if entry, ok := store.GetEntryInfo(r.URL.Path[len(uApiGet):]); !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "{}")
		} else {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(entry); err != nil {
				panic(err)
			}
		}
	})

	http.HandleFunc(uApiNew, func(w http.ResponseWriter, r *http.Request) {
		var entry StoreEntry
		log.Println(r)
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxData))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &entry); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		id := store.NewEntry(entry)
		log.Println(id)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(id); err != nil {
			panic(err)
		}
	})

	log.Fatal(http.ListenAndServe(listen, nil))
}

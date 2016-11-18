package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	schemeHost = "http://localhost"
	listen     = ":9154"
	uApiGet    = "/api/v1/get/"
	uApiNew    = "/api/v1/new"
	uGet       = "/g"
	uInfo      = "/i"
	uFav       = "/favicon.ico"
	maxData    = 1048576 // 1MB
)

func main() {
	store := make(secretStore)
	store.NewEntry("secret", 100, "test")

	/*
		/							# intro page
		/api/v1
			/api/v1/get/34g34g243	# get entry, decrease counter
			/api/v1/new				# generate new entry (takes POST data)
			/api/v1/info/34g34g243	# show info for entry
		/g?id=34g34g243				# show entry, decrease counter
		/n							# generate new entry (takes POST data)
		/i?id=34g34g243				# show info for entry (e. g. after creating)
	*/
	tView := template.New("view")
	tView.Parse(htmlMaster)
	tView.Parse(htmlView)
	tViewErr := template.New("viewErr")
	tViewErr.Parse(htmlMaster)
	tViewErr.Parse(htmlViewErr)
	tViewInfo := template.New("viewInfo")
	tViewInfo.Parse(htmlMaster)
	tViewInfo.Parse(htmlViewInfo)

	http.HandleFunc(uApiGet, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		id := r.URL.Path[len(uApiGet):]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if entry, ok := store.GetEntryInfo(id); !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "{}")
		} else {
			store.Click(id)
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
		id := store.AddEntry(entry, "")
		log.Println(id)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(id); err != nil {
			panic(err)
		}
	})

	http.HandleFunc(uGet, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		id := r.URL.Query().Get("id")
		if entry, ok := store.GetEntryInfo(id); !ok {
			w.WriteHeader(http.StatusNotFound)
			tViewErr.ExecuteTemplate(w, "master", nil)
		} else {
			store.Click(id)
			w.WriteHeader(http.StatusOK)
			tView.ExecuteTemplate(w, "master", entry)
		}
	})

	http.HandleFunc(uInfo, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		id := r.URL.Query().Get("id")
		if entry, ok := store.GetEntryInfo(id); !ok {
			w.WriteHeader(http.StatusNotFound)
			tViewErr.ExecuteTemplate(w, "master", nil)
		} else {
			w.WriteHeader(http.StatusOK)
			tViewInfo.ExecuteTemplate(w, "master", entry)
		}
	})

	http.HandleFunc(uFav, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(favicon)
	})

	log.Fatal(http.ListenAndServe(listen, nil))
}

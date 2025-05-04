package httpio

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sstark/gjfy/fileio"
	"github.com/sstark/gjfy/store"
)

type viewInfoEntry struct {
	store.StoreEntryInfo
	UserMessageView string
}

var htmlTemplates *template.Template

func init() {
	htmlTemplates, _ = template.ParseFS(fileio.HtmlTemplates, "*.tmpl")
}

func HandleIndex(fAllowAnonymous bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		type Data struct {
			AllowAnonymous bool
		}
		htmlTemplates.ExecuteTemplate(w, "index", &Data{AllowAnonymous: fAllowAnonymous})
	})
}

func HandleGet(memstore store.SecretStore, urlbase string, fNotify bool, userMessageView string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if entry, ok := memstore.GetEntryInfo(id, urlbase); !ok {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("entry not found: %s", id)
			htmlTemplates.ExecuteTemplate(w, "error", nil)
		} else {
			memstore.Click(id, r, fNotify)
			w.WriteHeader(http.StatusOK)
			viewEntry := viewInfoEntry{entry, userMessageView}
			htmlTemplates.ExecuteTemplate(w, "view", viewEntry)
		}
	})
}

func HandleInfo(memstore store.SecretStore, urlbase string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if entry, ok := memstore.GetEntryInfo(id, urlbase); !ok {
			w.WriteHeader(http.StatusNotFound)
			htmlTemplates.ExecuteTemplate(w, "error", nil)
		} else {
			w.WriteHeader(http.StatusOK)
			htmlTemplates.ExecuteTemplate(w, "info", entry)
		}
	})
}

func HandleCreate(memstore store.SecretStore, urlbase string, urlget string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entry := memstore.NewEntry(r.Form.Get("secret"), 1, 7, "anonymous", "")
		w.Write([]byte(fmt.Sprintf("%s%s?id=%s", urlbase, urlget, entry)))
	})
}

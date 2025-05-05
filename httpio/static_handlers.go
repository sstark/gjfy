package httpio

import (
	"bytes"
	"net/http"
	"time"

	"github.com/sstark/gjfy/fileio"
	"github.com/sstark/gjfy/misc"
)

func HandleStaticFav() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(fileio.Favicon)
	})
}

func HandleStaticLogoSmall() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(fileio.GjfyLogoSmall)
	})
}

func HandleStaticCss(cssp *[]byte, updatedp *time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		css := *cssp
		updated := *updatedp
		http.ServeContent(w, r, fileio.CssFileName, updated, bytes.NewReader(css))
	})
}

func HandleStaticLogo(logop *[]byte, updatedp *time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logo := *logop
		updated := *updatedp
		http.ServeContent(w, r, fileio.LogoFileName, updated, bytes.NewReader(logo))
	})
}

func HandleStaticClientShellScript(urlbase string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-sh")
		w.WriteHeader(http.StatusOK)
		misc.ClientShellScript(w, urlbase+ApiNew)
	})
}

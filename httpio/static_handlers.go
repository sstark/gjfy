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

func HandleStaticCss(css []byte, updated time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, fileio.CssFileName, updated, bytes.NewReader(css))
	})
}

func HandleStaticLogo(logo []byte, updated time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, fileio.LogoFileName, updated, bytes.NewReader(logo))
	})
}

func HandleStaticClientShellScript(urlbase, urlapinew string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-sh")
		w.WriteHeader(http.StatusOK)
		misc.ClientShellScript(w, urlbase+urlapinew)
	})
}

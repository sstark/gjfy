package misc

import "net/http"

func GetRealIP(r *http.Request) string {
	xFF := r.Header.Get("X-Forwarded-For")
	xRI := r.Header.Get("X-Real-IP")
	if xFF != "" {
		return xFF
	} else if xRI != "" {
		return xRI
	}
	return "none"
}

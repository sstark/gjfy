package misc

import (
	"net/http"
	"testing"
)

type GetRealIPTestCase struct {
	request http.Request
	headers map[string]string
	expected string
}

func makeFauxRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return req
}

func TestGetRealIP(t *testing.T) {
	tests := map[string]GetRealIPTestCase{
		"no proxy headers": {
			*makeFauxRequest("/api/v1/foo"),
			make(map[string]string),
			"none",
		},
		"forwarded-for": {
			*makeFauxRequest("/api/v1/fie"),
			map[string]string{
				"X-Forwarded-For": "123.123.123.123",
			},
			"123.123.123.123",
		},
		"real-ip": {
			*makeFauxRequest("/api/v1/bar"),
			map[string]string{
				"X-Real-IP": "222.222.111.111",
			},
			"222.222.111.111",
		},
		"both": {
			*makeFauxRequest("/api/v1/baz"),
			map[string]string{
				"X-Real-IP": "101.101.202.202",
				"X-Forwarded-For": "144.144.133.133",
			},
			"144.144.133.133",
		},
	}
	for label, test := range tests {
		t.Run(label, func(t *testing.T) {
			for header, val := range test.headers {
				test.request.Header.Add(header, val)
			}
			got := GetRealIP(&test.request)
			if got != test.expected {
				t.Errorf("Got a different header value (%s) than expected (%s)", got, test.expected)
			}
		})
	}
}

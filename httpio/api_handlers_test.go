package httpio

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/sstark/gjfy/store"
	"github.com/sstark/gjfy/tokendb"
)

var mockNow = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func TestHandleApiGet(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	secretStore := make(store.SecretStore)
	secretStore.NewEntry("secret", 1, 1, "auth", "testid")
	urlbase := "http://localhost:9154"
	req, _ := http.NewRequest("GET", urlbase+ApiGet+"testid", nil)
	rr := httptest.NewRecorder()
	handler := HandleApiGet(secretStore, urlbase, false)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusOK)
	}
	expected := `{"secret":"secret","max_clicks":1,"clicks":0,"date_added":"2009-11-10T23:00:00Z","valid_for":1,"auth_token":"auth","id":"testid","path_query":"/g?id=testid","url":"http://localhost:9154/g?id=testid","api_url":"http://localhost:9154/api/v1/get/testid"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHandleApiGetNonExisting(t *testing.T) {
	secretStore := make(store.SecretStore)
	urlbase := "http://localhost:9154"
	req, _ := http.NewRequest("GET", urlbase+ApiGet+"foo", nil)
	rr := httptest.NewRecorder()
	handler := HandleApiGet(secretStore, urlbase, false)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusOK)
	}
	expected := `{"error":"not found"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHandleApiNew(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	secretStore := make(store.SecretStore)
	auth := tokendb.MakeTokenDB([]byte(`[{
					"token": "footoken",
					"email": "test@example.org"
				}]`))
	urlbase := "http://localhost:9154"
	postdata := bytes.NewReader([]byte(`{
					"auth_token": "footoken",
					"secret": "sekrit",
					"max_clicks": 3 
				}`))
	req, _ := http.NewRequest("POST", urlbase+ApiNew, postdata)
	rr := httptest.NewRecorder()
	handler := HandleApiNew(secretStore, urlbase, &auth)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusCreated)
	}
	expected := `{"secret":"#HIDDEN#","max_clicks":3,"clicks":0,"date_added":"2009-11-10T23:00:00Z","valid_for":7,"auth_token":"test@example.org","id":"EUwXrkDvd1Gw2jNG-gvRr68rGaaNIeJoJOpLQ2WTqNI","path_query":"/g?id=EUwXrkDvd1Gw2jNG-gvRr68rGaaNIeJoJOpLQ2WTqNI","url":"http://localhost:9154/g?id=EUwXrkDvd1Gw2jNG-gvRr68rGaaNIeJoJOpLQ2WTqNI","api_url":"http://localhost:9154/api/v1/get/EUwXrkDvd1Gw2jNG-gvRr68rGaaNIeJoJOpLQ2WTqNI"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHandleApiNewUnauthorized(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	secretStore := make(store.SecretStore)
	auth := tokendb.MakeTokenDB([]byte(`[{
					"token": "footoken",
					"email": "test@example.org"
				}]`))
	urlbase := "http://localhost:9154"
	postdata := bytes.NewReader([]byte(`{
					"auth_token": "wrongtoken",
					"secret": "sekrit",
					"max_clicks": 3
				}`))
	req, _ := http.NewRequest("POST", urlbase+ApiNew, postdata)
	rr := httptest.NewRecorder()
	handler := HandleApiNew(secretStore, urlbase, &auth)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusUnauthorized)
	}
	expected := `{"error":"unauthorized"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

type mockFailingResponseWriter struct {
	statusCode int
	headers    http.Header
}

func (m *mockFailingResponseWriter) WriteHeader(_ int) {
	return
}

func (m *mockFailingResponseWriter) Header() http.Header {
	if m.headers == nil {
		m.headers = make(http.Header)
	}
	return m.headers
}

func (m *mockFailingResponseWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("simulated write error")
}

func TestHandleApiNewMalformed(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	secretStore := make(store.SecretStore)
	auth := tokendb.MakeTokenDB([]byte(`[{
					"token": "footoken",
					"email": "test@example.org"
				}]`))
	urlbase := "http://localhost:9154"
	postdata := bytes.NewReader([]byte(`{
					"auth_token": "wrongtoken",
					"secret": 24,
					"max_clicks": "baz"
				}`))
	req, _ := http.NewRequest("POST", urlbase+ApiNew, postdata)
	rr := httptest.NewRecorder()
	handler := HandleApiNew(secretStore, urlbase, &auth)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusUnprocessableEntity)
	}
	expected := `{"error":"json: cannot unmarshal number into Go struct field StoreEntry.secret of type string"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestJsonRespond(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		rr := httptest.NewRecorder()
		type testContent struct {
			SomeValue string `json:"somevalue"`
		}

		jsonRespond(rr, http.StatusOK, testContent{"foobar"})

		expectedContentType := "application/json; charset=UTF-8"
		if rr.Header().Get("Content-Type") != expectedContentType {
			t.Errorf("handler returned wrong content type: got %v, wanted %v", rr.Header().Get("Content-Type"), expectedContentType)
		}
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v, wanted %v", rr.Code, http.StatusOK)
		}
		expectedBody := `{"somevalue":"foobar"}
`
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got\n%v want\n%v", rr.Body.String(), expectedBody)
		}
	})

	t.Run("unhappy case: json encoding error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		var invalidData chan int

		jsonRespond(rr, http.StatusOK, invalidData)

		expectedBody := `{"error":"internal error"}`
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got\n%v want\n%v", rr.Body.String(), expectedBody)
		}
	})

	t.Run("unhappy case: can not write", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()
		rr := &mockFailingResponseWriter{statusCode: http.StatusOK, headers: make(http.Header)}

		jsonRespond(rr, http.StatusOK, "foo")

		expectedError := `error writing response: simulated write error
`
		if !strings.HasSuffix(buf.String(), expectedError) {
			t.Errorf("handler did not log error: got\n%v wanted suffix\n%v", buf.String(), expectedError)
		}
	})
}

package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/sstark/gjfy/fileio"
	"github.com/sstark/gjfy/httpio"
	"github.com/sstark/gjfy/store"
	"github.com/sstark/gjfy/tokendb"
)

var mockNow = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func createAuthDBFile(dir, content string) {
	authDBFile, _ := os.Create(filepath.Join(dir, tokendb.AuthFileName))
	defer authDBFile.Close()
	authDBFile.WriteString(content)
}

func createUserMessageViewFile(dir, content string) {
	umvFile, _ := os.Create(filepath.Join(dir, fileio.UserMessageViewFilename))
	defer umvFile.Close()
	umvFile.WriteString(content)
}

func TestUpdateFiles(t *testing.T) {
	tmpdir, _ := os.MkdirTemp("", "gjfy_test")
	os.Chdir(tmpdir)
	createAuthDBFile(tmpdir, `[{
        "token": "test",
        "email": "test@example.org"
    },
    {
        "token": "test2",
        "email": "other@example.org"
    }
]`)
	auth1, css1, logo1, userMessageView1, updated1 := updateFiles()
	time.Sleep(time.Millisecond * 2)
	auth2, css2, logo2, userMessageView2, updated2 := updateFiles()
	if !reflect.DeepEqual(auth1, auth2) || !reflect.DeepEqual(css1, css2) || !reflect.DeepEqual(logo1, logo2) || userMessageView1 != userMessageView2 {
		t.Errorf("Running updateFiles twice gives differing results")
	}
	if updated1 == updated2 {
		t.Errorf("Timestamp did not change between updates")
	}
	createAuthDBFile(tmpdir, `[{
        "token": "test",
        "email": "test@example.org"
    },
    {
        "token": "test3",
        "email": "foobar@example.org"
    }
]`)
	auth3, _, _, _, _ := updateFiles()
	if reflect.DeepEqual(auth2, auth3) {
		t.Errorf("auth.db was not updated after changing file")
	}
	userMessage := "foo bar baz!"
	createUserMessageViewFile(tmpdir, userMessage)
	_, _, _, userMessageView3, _ := updateFiles()
	if userMessageView3 != userMessage {
		t.Errorf("userMessageView was not updated from file")
	}
}

// Try posting a secret without proper token, then update the auth.db
// file with the correct token, reload it and make sure posting works.
func TestAuthDBUpdatedAtRuntime(t *testing.T) {
	tmpdir, _ := os.MkdirTemp("", "gjfy_test")
	os.Chdir(tmpdir)
	createAuthDBFile(tmpdir, `[{
        "token": "test",
        "email": "test@example.org"
    }
]`)
	auth, _, _, _, _ := updateFiles()
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	store := make(store.SecretStore)
	urlbase := "http://localhost:9154"
	postdata := bytes.NewReader([]byte(`{
					"auth_token": "sometoken",
					"secret": "sekrit",
					"max_clicks": 3 
				}`))
	req, _ := http.NewRequest("POST", urlbase+httpio.ApiNew, postdata)
	rr := httptest.NewRecorder()
	handler := httpio.HandleApiNew(store, urlbase, &auth)
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
	createAuthDBFile(tmpdir, `[{
        "token": "test",
        "email": "test@example.org"
    },
    {
        "token": "sometoken",
        "email": "other@example.org"
    }
]`)
	auth, _, _, _, _ = updateFiles()
	postdata2 := bytes.NewReader([]byte(`{
					"auth_token": "sometoken",
					"secret": "sekrit",
					"max_clicks": 3 
				}`))
	req2, _ := http.NewRequest("POST", urlbase+httpio.ApiNew, postdata2)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v, wanted %v", status, http.StatusCreated)
	}
	expected2 := `{"secret":"#HIDDEN#","max_clicks":3,"clicks":0,"date_added":"2009-11-10T23:00:00Z","valid_for":7,"auth_token":"other@example.org","id":"5SxjhlIQghCDo4pVktIqVKsti1TGqJ5O9g6eEJMTyhA","path_query":"/g?id=5SxjhlIQghCDo4pVktIqVKsti1TGqJ5O9g6eEJMTyhA","url":"http://localhost:9154/g?id=5SxjhlIQghCDo4pVktIqVKsti1TGqJ5O9g6eEJMTyhA","api_url":"http://localhost:9154/api/v1/get/5SxjhlIQghCDo4pVktIqVKsti1TGqJ5O9g6eEJMTyhA"}
`
	if rr2.Body.String() != expected2 {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v", rr2.Body.String(), expected2)
	}
}

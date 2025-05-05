package tokendb

import (
	"io"
	"log"
	"testing"

	"github.com/sstark/gjfy/store"
)

type tdbTestPair struct {
	in  []byte
	out bool
}

var tdbTestPairs = []tdbTestPair{
	{
		in:  []byte("bla"),
		out: false,
	},
	{
		in: []byte(`[{
				"token": "test",
				"email": "test@example.org"
			},
			{
				"token": "test2",
				"email": "other@example.org"
			}]`),
		out: true,
	},
}

func TestAuth_makeTokenDB(t *testing.T) {
	log.SetOutput(io.Discard)
	for _, pair := range tdbTestPairs {
		tdb := MakeTokenDB(pair.in)
		if (tdb != nil) != pair.out {
			t.Errorf("%v should be %v", tdb, pair.out)
		}
	}
}

type tdbFindTokenTest struct {
	db  []byte
	in  string
	out string
}

func TestAuthFindToken(t *testing.T) {
	tdbFindTokenTests := map[string]tdbFindTokenTest{
		"look for non-existing entry": {
			db:  []byte("bla"),
			in:  "foo@example.com",
			out: "",
		},
		"look for existing entry": {
			db: []byte(`[{
					"token": "test",
					"email": "test@example.org"
				},
				{
					"token": "test2",
					"email": "other@example.org"
				}]`),
			in:  "test2",
			out: "other@example.org",
		},
	}
	for label, test := range tdbFindTokenTests {
		t.Run(label, func(t *testing.T) {
			tdb := MakeTokenDB(test.db)
			out := tdb.findToken(test.in)
			if out != test.out {
				t.Errorf("unexpected output: expected %v, got %v", test.out, out)
			}
		})
	}
}

type tdbIsAuthorizedTest struct {
	db  []byte
	in  string
	found bool
	email string
	out bool
}

func TestIsAuthorized(t *testing.T) {
	tdbIsAuthorizedTests := map[string]tdbIsAuthorizedTest{
		"look for non-existing entry": {
			db:  []byte("bla"),
			in:  "foobar",
			found: false,
			email: "",
			out: false,
		},
		"look for existing entry": {
			db: []byte(`[{
					"token": "test",
					"email": "test@example.org"
				},
				{
					"token": "test2",
					"email": "other@example.org"
				}]`),
			in:  "id",
			found: true,
			email: "other@example.org",
			out: true,
		},
	}
	store := make(store.SecretStore)
	store.NewEntry("secret", 1, 1, "test2", "id")
	for label, test := range tdbIsAuthorizedTests {
		t.Run(label, func(t *testing.T) {
			entry, ok := store.GetEntry(test.in)
			if ok != test.found {
				t.Errorf("when finding store entry: expected %v, got %v", test.found, ok)
			}
			tdb := MakeTokenDB(test.db)
			out := tdb.IsAuthorized(&entry)
			if out != test.out {
				t.Errorf("unexpected output: expected %v, got %v", test.out, out)
			}
			// A side-effect of isAuthorized is to change the token into the email address, check for that
			if ok {
				if entry.AuthToken != test.email {
					t.Errorf("entry AuthToken field was not changed to expected value: expected %v, got %v", test.email, entry.AuthToken)
				}
			}
		})
	}
}

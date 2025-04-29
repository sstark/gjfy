package tokendb

import (
	"io"
	"log"
	"testing"
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

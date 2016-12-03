package main

import (
	"github.com/bouk/monkey"
	"reflect"
	"testing"
	"time"
)

var mockNow = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func TestHashStruct(t *testing.T) {
	var in = StoreEntry{"secret1", 5, 0, mockNow, 3, "authtoken"}
	var wanted = "0Y3Mkcz36xM0hwrSnVw3PMebEMfa27Oi1mmuaELD4-Q"
	got := hashStruct(in)
	if got != wanted {
		t.Errorf("got %v, wanted %v", got, wanted)
	}
}

type StoreEntryInput struct {
	secret    string
	maxClicks int
	validFor  int
	authToken string
	id        string
}

type StoreEntryOutput struct {
	StoreEntry
	id string
}

type StoreEntryTestPair struct {
	in  StoreEntryInput
	out StoreEntryOutput
}

var StoreEntryTestPairs = []StoreEntryTestPair{
	{
		StoreEntryInput{"secret1", 5, 3, "authtoken", "id1"},
		StoreEntryOutput{StoreEntry{"secret1", 5, 0, mockNow, 3, "authtoken"}, "id1"},
	},
	{
		StoreEntryInput{"secret2", 2, 3, "authtoken", ""},
		StoreEntryOutput{StoreEntry{"secret2", 2, 0, mockNow, 3, "authtoken"}, "iLbLBYFzULLUfB84p8VHldWd4VnHg0mZq_5S45p0lEk"},
	},
}

func TestStore_NewEntry(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return mockNow
	})
	defer monkey.Unpatch(time.Now)
	store := make(secretStore)
	for _, p := range StoreEntryTestPairs {
		outId := store.NewEntry(p.in.secret, p.in.maxClicks, p.in.validFor, p.in.authToken, p.in.id)
		if outId != p.out.id {
			t.Errorf("got %v, wanted %v", outId, p.out.id)
		}
		outEntry, ok := store.GetEntry(outId)
		if !ok {
			t.Errorf("new entry not found under %v", outId)
		}
		if !reflect.DeepEqual(p.out.StoreEntry, outEntry) {
			t.Errorf("got %v, wanted %v", p.out.StoreEntry, outEntry)
		}
	}
}

func TestStore_GetEntryInfo(t *testing.T) {
	store := make(secretStore)
	store.NewEntry("secret", 1, 1, "auth", "testid")
	out, ok := store.GetEntryInfo("testid")
	if !ok {
		t.Errorf("new entry not found under %v", "testid")
	}
	wanted := "http://localhost:/api/v1/get/testid"
	if out.ApiUrl != wanted {
		t.Errorf("got %v, wanted %v", out.ApiUrl, wanted)
	}
	wanted = "http://localhost:/g?id=testid"
	if out.Url != wanted {
		t.Errorf("got %v, wanted %v", out.Url, wanted)
	}
}

package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/sstark/gjfy/fileio"
	"github.com/sstark/gjfy/tokendb"
)

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

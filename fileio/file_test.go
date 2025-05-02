package fileio

import (
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestFileOrFunc(t *testing.T) {
	log.SetOutput(io.Discard)
	byteval := []byte{0x1, 0x2, 0x3, 0x4}
	fn := "/this/file/does/not/exist/anywhere.xxxxx"
	bytes_returned := FileOrFunc(fn, func(fn string) []byte {
		return byteval
	})
	if !reflect.DeepEqual(bytes_returned, byteval) {
		t.Errorf("FileOrFunc did not fall back properly to the output of the callback function. Got: %v, Expected: %v", bytes_returned, byteval)
	}
}

func TestTryReadFile(t *testing.T) {
	tmpdir, _ := os.MkdirTemp("", "gjfy_test")
	defer os.Remove(tmpdir)

	fileName := "gjfy_testfile"
	testContent := []byte("test")
	testFileA := path.Join(tmpdir, fileName)
	os.WriteFile(testFileA, testContent, 0644)
	defer os.Remove(testFileA)

	configDir = tmpdir
	log.SetOutput(io.Discard)

	// Test file in configDir
	bytes := TryReadFile(fileName)
	if string(bytes) != string(testContent) {
		t.Errorf("got %v, wanted %v", bytes, testContent)
	}

	// Test file does not exist
	bytes = TryReadFile("doesnotexist")
	if string(bytes) != "" {
		t.Errorf("got %v, wanted empty slice", bytes)
	}

	// Test file in pwd shadows file in configDir
	testContentPwd := []byte("testPwd")
	wd, _ := os.Getwd()
	testFileB := path.Join(wd, fileName)
	os.WriteFile(testFileB, testContentPwd, 0644)
	defer os.Remove(testFileB)
	bytes = TryReadFile(fileName)
	if string(bytes) != string(testContentPwd) {
		t.Errorf("got %v, wanted %v", bytes, testContentPwd)
	}
}

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestTryReadFile(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "gjfy_test")
	defer os.Remove(tmpdir)

	fileName := "gjfy_testfile"
	testContent := []byte("test")
	testFileA := path.Join(tmpdir, fileName)
	ioutil.WriteFile(testFileA, testContent, 0644)
	defer os.Remove(testFileA)

	configDir = tmpdir
	log.SetOutput(ioutil.Discard)

	// Test file in configDir
	bytes := tryReadFile(fileName)
	if string(bytes) != string(testContent) {
		t.Errorf("got %v, wanted %v", bytes, testContent)
	}

	// Test file does not exist
	bytes = tryReadFile("doesnotexist")
	if string(bytes) != "" {
		t.Errorf("got %v, wanted empty slice", bytes)
	}

	// Test file in pwd shadows file in configDir
	testContentPwd := []byte("testPwd")
	wd, _ := os.Getwd()
	testFileB := path.Join(wd, fileName)
	ioutil.WriteFile(testFileB, testContentPwd, 0644)
	defer os.Remove(testFileB)
	bytes = tryReadFile(fileName)
	if string(bytes) != string(testContentPwd) {
		t.Errorf("got %v, wanted %v", bytes, testContentPwd)
	}
}

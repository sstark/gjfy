package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

// tryReadFile takes a _filename_ as an argument and tries to read that file
// from cwd and then configDir. If it doesn't find the file in any of them it
// will return an empty byte slice.
func tryReadFile(fn string) []byte {
	var dirs []string
	cwd, err := os.Getwd()
	if err == nil {
		dirs = append(dirs, cwd)
	} else {
		log.Println("could not get working directory")
	}
	dirs = append(dirs, configDir)
	for _, dir := range dirs {
		contents, err := ioutil.ReadFile(path.Join(dir, fn))
		if err == nil {
			log.Printf("found %s in %s\n", fn, dir)
			return contents
		}
	}
	log.Println("could not find", fn)
	return []byte{}
}

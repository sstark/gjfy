package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

func NotifyMail(to, msg string) {
	if fNotify {
		go SendMail(to, "GJFY notice", msg)
	}
}

func SendMail(to, subject, msg string) {
	sendmail := exec.Command("mail", "-s", subject, to)
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}
	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	sendmail.Start()
	stdin.Write([]byte(msg))
	stdin.Write([]byte("\n"))
	stdin.Close()
	ioutil.ReadAll(stdout)
	sendmail.Wait()
	log.Printf("sending notification to %s done.\n", to)
}

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func NotifyMail(to, msg string) {
	SendMail(to, "GJFY notice", msg)
}

func SendMail(to, subject, msg string) {
	sendmail := exec.Command("mail", "-s", subject, to)
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		panic(err)
	}
	sendmail.Start()
	stdin.Write([]byte(msg))
	stdin.Write([]byte("\n"))
	stdin.Close()
	sentBytes, _ := ioutil.ReadAll(stdout)
	sendmail.Wait()
	fmt.Println("Send Command Output\n")
	fmt.Println(string(sentBytes))
}

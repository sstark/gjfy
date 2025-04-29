package misc

import (
	"io"
	"log"
	"os/exec"
)

func NotifyMail(to, msg string) {
	go SendMail(to, "GJFY notice", msg)
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
	io.ReadAll(stdout)
	sendmail.Wait()
	log.Printf("sending notification to %s done.\n", to)
}

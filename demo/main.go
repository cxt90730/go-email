package main

import "github.com/cxt90730/go-email"

func main() {
	s, err := email_sdk.NewEmailService("mail.conf", "Email")
	if err != nil {
		panic(err)
	}
	msg := s.NewMessage("cxt", "test", "test content", []string{"cxt@onebooktech.com"}, []string{"cxt"})
	err = s.SendMail(msg, "")
	if err != nil {
		panic(err)
	}
}

package email_sdk

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/robfig/config"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"
)

type EmailService struct {
	Addr     string
	Port     string
	AuthName string
	AuthPwd  string
}

type MailMessage struct {
	SenderName string
	Sender     string
	To         []string
	ToName     []string
	Subject    string
	Body       string
	Marker     string
}

func NewEmailService(confPath string, tag string) (*EmailService, error) {
	c, err := config.ReadDefault(confPath)
	if err != nil {
		return nil, err
	}
	emailService := new(EmailService)
	emailService.Addr, err = c.String(tag, "addr")
	if err != nil {
		return nil, err
	}
	emailService.Port, err = c.String(tag, "port")
	if err != nil {
		return nil, err
	}
	emailService.AuthName, err = c.String(tag, "auth_name")
	if err != nil {
		return nil, err
	}
	emailService.AuthPwd, err = c.String(tag, "auth_pwd")
	if err != nil {
		return nil, err
	}
	return emailService, nil
}

// mail headers
func (m *MailMessage) DefaultHead() string {
	return fmt.Sprintf("From: %s <%s>\r\nTo: %s <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n--%s",
		m.SenderName, m.Sender, m.ToName[0], m.To[0], m.Subject, m.Marker, m.Marker)
}

// body (text or HTML)
func (m *MailMessage) DefaultBodys() string {
	return fmt.Sprintf("\r\nContent-Type: text/html\r\nContent-Transfer-Encoding:8bit\r\n\r\n%s\r\n--%s", m.Body, m.Marker)
}

var ContentType = map[string]string{
	".gif":  "image/gif",
	".doc":  "application/msword",
	".docx": "application/msword",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	//ELSE?
}

func (m *EmailService) NewMessage(name, subject, body string, to, ToName []string) *MailMessage {
	return &MailMessage{
		SenderName: name,
		Sender:     m.AuthName,
		To:         to,
		ToName:     ToName,
		Subject:    subject,
		Body:       body,
		Marker:     "Onebooktech",
	}
}

//send file if you need.
//otherwise filename is empty.

func (m *MailMessage) Encode(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}
	var buf bytes.Buffer

	name := filepath.Base(filename)
	contentType := ContentType[strings.ToLower(filepath.Ext(filename))]

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	lineMaxLength := 500
	nbrLines := len(encoded) / lineMaxLength

	//append lines to buffer
	for i := 0; i < nbrLines; i++ {
		buf.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")
	}

	//append last line in buffer
	buf.WriteString(encoded[nbrLines*lineMaxLength:])

	return fmt.Sprintf("\r\nContent-Type: %s; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n--%s--",
		contentType, name, name, buf.String(), m.Marker), nil
}

func (m *EmailService) SendMail(msg *MailMessage, fileName string) error {
	encodeStr, err := msg.Encode(fileName)
	if err != nil {
		return err
	}

	content := msg.DefaultHead() + msg.DefaultBodys() + encodeStr

	auth := smtp.PlainAuth(msg.Sender, m.AuthName, m.AuthPwd, m.Addr)

	host := m.Addr + ":" + m.Port
	if err := smtp.SendMail(host, auth, msg.Sender, msg.To, []byte(content)); err != nil {
		return err
	}
	return nil
}

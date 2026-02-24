package notifications

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"text/template"
	"time"
)

type EmailSender struct {
	Addr     string
	Username string
	Password string
	From     string
	Auth     smtp.Auth
}

func NewEmailSender(addr, username, password, from string) *EmailSender {
	host := addr
	if i := bytes.IndexByte([]byte(addr), ':'); i > 0 {
		host = addr[:i]
	}
	return &EmailSender{Addr: addr, Username: username, Password: password, From: from, Auth: smtp.PlainAuth("", username, password, host)}
}

func (s *EmailSender) SendTemplate(ctx context.Context, to, subject, tpl string, data any, retries int) error {
	t, err := template.New("mail").Parse(tpl)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", to, subject, body.String())
	var sendErr error
	for i := 0; i <= retries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		sendErr = smtp.SendMail(s.Addr, s.Auth, s.From, []string{to}, []byte(msg))
		if sendErr == nil {
			return nil
		}
		time.Sleep(time.Duration(1<<i) * time.Second)
	}
	return sendErr
}

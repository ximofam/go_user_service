package mail

import (
	"context"

	"gopkg.in/gomail.v2"
)

type MailtrapSender struct {
	dialer *gomail.Dialer
}

func NewMailtrapSender(host string, port int, username, password string) *MailtrapSender {
	return &MailtrapSender{
		dialer: gomail.NewDialer(host, port, username, password),
	}
}

func (s *MailtrapSender) Send(ctx context.Context, mail *Mail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", mail.Subject)
	if mail.TextBody != "" {
		m.SetBody("text/plain", mail.TextBody)
	}
	if mail.HTMLBody != "" {
		if mail.TextBody != "" {
			m.AddAlternative("text/html", mail.HTMLBody)
		} else {
			m.SetBody("text/html", mail.HTMLBody)
		}
	}

	return s.dialer.DialAndSend(m)
}

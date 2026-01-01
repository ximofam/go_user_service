package mail

import "context"

type MailService interface {
	Send(ctx context.Context, mail *Mail) error
}

type Mail struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	TextBody string   `json:"text_body"`
	HTMLBody string   `json:"html_body"`
}

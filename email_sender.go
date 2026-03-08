package emailnotification

import "context"

type EmailSendInput struct {
	From     string
	FromName string
	To       []string
	Cc       []string
	Bcc      []string
	ReplyTo  string
	Subject  string
	HTMLBody string
	TextBody string
}

type EmailSender interface {
	Send(ctx context.Context, input EmailSendInput) error
}

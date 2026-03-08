package sender

import (
	"context"

	emailnotification "github.com/KOMKZ/go-yogan-domain-email-notification"
	email "github.com/KOMKZ/go-yogan-component-email"
)

type ComponentEmailSender struct {
	manager *email.Manager
}

func NewComponentEmailSender(manager *email.Manager) *ComponentEmailSender {
	return &ComponentEmailSender{manager: manager}
}

func (s *ComponentEmailSender) Send(ctx context.Context, input emailnotification.EmailSendInput) error {
	builder := s.manager.New()

	for _, to := range input.To {
		builder.To(to)
	}
	builder.Subject(input.Subject).Body(input.HTMLBody)

	if input.TextBody != "" {
		builder.BodyText(input.TextBody)
	}
	if input.From != "" {
		builder.From(input.From)
	}
	if input.FromName != "" {
		builder.FromName(input.FromName)
	}
	for _, cc := range input.Cc {
		builder.Cc(cc)
	}
	for _, bcc := range input.Bcc {
		builder.Bcc(bcc)
	}
	if input.ReplyTo != "" {
		builder.ReplyTo(input.ReplyTo)
	}

	_, err := builder.Send(ctx)
	return err
}

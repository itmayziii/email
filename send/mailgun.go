package send

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSenderAdapter allows a mailgun.Mailgun interface to become compatible with the Sender interface.
type MailgunSenderAdapter struct {
	mailgun mailgun.Mailgun
}

func (adapter MailgunSenderAdapter) Send(ctx context.Context, m Message) (string, error) {
	message := adapter.mailgun.NewMessage(m.Sender, m.Subject, "", m.To...)
	message.SetHtml(m.Body)

	for _, cc := range m.Cc {
		message.AddCC(cc)
	}

	for _, bcc := range m.Bcc {
		message.AddCC(bcc)
	}

	_, id, err := adapter.mailgun.Send(ctx, message)
	return id, err
}

// NewMailgunSender constructs a MailgunSenderAdapter
func NewMailgunSender(mailgun mailgun.Mailgun) Sender {
	return MailgunSenderAdapter{mailgun: mailgun}
}

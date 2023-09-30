package send

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSenderAdapter allows a mailgun.Mailgun interface to become compatible with the Sender interface.
type MailgunSenderAdapter struct {
	mailgun mailgun.Mailgun
}

func (adapter MailgunSenderAdapter) Send(ctx context.Context, m Message) (string, string, error) {
	message := adapter.mailgun.NewMessage(m.Sender, m.Subject, "", m.To...)
	message.SetHtml(m.Body)
	return adapter.mailgun.Send(ctx, message)
}

// NewMailgunSender constructs a MailgunSenderAdapter
func NewMailgunSender(mailgun mailgun.Mailgun) Sender {
	return MailgunSenderAdapter{mailgun: mailgun}
}

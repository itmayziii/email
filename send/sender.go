package send

import "context"

// Sender sends an email with the provided Message.
type Sender interface {
	Send(ctx context.Context, m Message) (string, string, error)
}

// Message represents an email.
type Message struct {
	Sender  string
	Subject string
	Body    string
	To      []string
}
